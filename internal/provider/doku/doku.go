package doku

import (
	"bytes"
	"context"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pandudpn/go-payment-gateway/internal/utils"
	"github.com/pandudpn/go-payment-gateway"
)

const (
	// ProviderName is the name of the Doku provider
	ProviderName = "doku"

	// API endpoints
	generatePaymentUri = "/payments/v2"
	statusUri          = "/transactions/v2"
	tokenUri           = "/authorization/v1/access-token/b2b"

	// header names
	headerSignature = "Signature"
	headerTimestamp = "Request-Timestamp"
	headerClientID  = "Client-Id"
	headerRequestID = "Request-Id"
	headerXClientKey = "X-CLIENT-KEY"
	headerXTimestamp = "X-TIMESTAMP"
	headerXSignature = "X-SIGNATURE"

	// API URLs
	sandboxURL    = "https://api-sandbox.doku.com"
	productionURL = "https://api.doku.com"
)

func init() {
	// Register this provider with the pg package
	pg.RegisterProvider(ProviderName, New)
}

type doku struct {
	config  *pg.ProviderConfig
	mapper  *Mapper
	httpCli *http.Client
}

// New creates a new Doku provider
func New(cfg *pg.ProviderConfig) (pg.Provider, error) {
	if cfg.ServerKey == "" || cfg.ClientKey == "" {
		return nil, pg.ErrMissingCredentials
	}

	return &doku{
		config:  cfg,
		mapper:  &Mapper{},
		httpCli: &http.Client{Timeout: getTimeout(cfg)},
	}, nil
}

// Name returns the provider name
func (d *doku) Name() string {
	return ProviderName
}

// CreateCharge creates a new payment transaction
func (d *doku) CreateCharge(ctx context.Context, params pg.ChargeParams) (*pg.ChargeResponse, error) {
	// Validate parameters
	if err := d.validateChargeParams(&params); err != nil {
		return nil, err
	}

	req := d.mapper.mapToGenerateRequest(params)
	responseBody, err := d.generatePayment(ctx, req)
	if err != nil {
		return nil, err
	}

	var resp GeneratePaymentResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.ResponseCode != "00" && resp.ResponseCode != "200" {
		return nil, pg.WrapProviderError(ProviderName, resp.ResponseCode, resp.ResponseMessage, nil)
	}

	return d.mapper.mapToChargeResponse(&resp, params.PaymentType), nil
}

// generatePayment initiates a payment
func (d *doku) generatePayment(ctx context.Context, params *GeneratePaymentRequest) ([]byte, error) {
	baseURL := d.getBaseURL()
	fullURL := baseURL + generatePaymentUri

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Generate ISO8601 timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)
	// Generate unique request ID
	requestID := generateRequestID()

	digest := d.generateDigest(bodyBytes)
	signature := d.generateSignature(digest, timestamp, requestID, generatePaymentUri)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerClientID, d.config.ClientKey)
	req.Header.Set(headerRequestID, requestID)
	req.Header.Set(headerTimestamp, timestamp)
	req.Header.Set(headerSignature, signature)

	resp, err := d.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// GetStatus retrieves payment status
func (d *doku) GetStatus(ctx context.Context, orderID string) (*pg.PaymentStatus, error) {
	baseURL := d.getBaseURL()
	fullURL := baseURL + statusUri

	req := &TransactionStatusRequest{
		TransactionID: orderID,
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Generate ISO8601 timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)
	// Generate unique request ID
	requestID := generateRequestID()

	digest := d.generateDigest(bodyBytes)
	signature := d.generateSignature(digest, timestamp, requestID, statusUri)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set(headerClientID, d.config.ClientKey)
	httpReq.Header.Set(headerRequestID, requestID)
	httpReq.Header.Set(headerTimestamp, timestamp)
	httpReq.Header.Set(headerSignature, signature)

	resp, err := d.httpCli.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(responseBody))
	}

	var dokuResponse TransactionStatusResponse
	if err := json.Unmarshal(responseBody, &dokuResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return d.mapper.mapToPaymentStatus(orderID, &dokuResponse), nil
}

// Cancel cancels a transaction
func (d *doku) Cancel(ctx context.Context, orderID string) error {
	// Doku doesn't have a direct cancel API
	// The payment will expire automatically based on VA expiration
	return fmt.Errorf("cancel not supported by Doku, payment will expire automatically")
}

// VerifyWebhook verifies webhook signature
func (d *doku) VerifyWebhook(r *http.Request) bool {
	signature := r.Header.Get("Signature")
	timestamp := r.Header.Get("Request-Timestamp")
	requestID := r.Header.Get("Request-Id")

	if signature == "" || timestamp == "" || requestID == "" {
		return false
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// Restore body for subsequent reads
	r.Body = io.NopCloser(bytes.NewReader(body))

	// Calculate digest
	digest := d.generateDigest(body)

	// Calculate expected signature
	expectedSignature := d.generateSignature(digest, timestamp, requestID, "/payments/v2")

	// Compare signatures
	return signature == expectedSignature
}

// ParseWebhook parses webhook payload
func (d *doku) ParseWebhook(r *http.Request) (*pg.WebhookEvent, error) {
	var webhookData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
		return nil, pg.ErrInvalidPayload
	}

	orderID, _ := webhookData["transaction_id"].(string)
	status, _ := webhookData["transaction_status"].(string)

	var amount int64
	if amountFloat, ok := webhookData["amount"].(float64); ok {
		amount = int64(amountFloat)
	} else if amountStr, ok := webhookData["amount"].(string); ok {
		amount, _ = strconv.ParseInt(amountStr, 10, 64)
	}

	// Parse timestamp
	var timestamp time.Time
	if paymentDate, ok := webhookData["payment_date_time"].(string); ok {
		timestamp, _ = time.Parse(time.RFC3339, paymentDate)
	}

	// Map status
	var mappedStatus pg.Status
	switch PaymentStatus(status) {
	case StatusSuccess:
		mappedStatus = pg.StatusSuccess
	case StatusPending:
		mappedStatus = pg.StatusPending
	case StatusFailed:
		mappedStatus = pg.StatusFailed
	case StatusCancelled:
		mappedStatus = pg.StatusCancelled
	default:
		mappedStatus = pg.StatusPending
	}

	return &pg.WebhookEvent{
		OrderID:       orderID,
		TransactionID: orderID,
		Status:        mappedStatus,
		Amount:        amount,
		EventType:     d.mapper.mapEventType(PaymentStatus(status)),
		Timestamp:     timestamp,
		Raw:           webhookData,
	}, nil
}

// GetToken retrieves an OAuth Bearer token using asymmetric RSA signature
// This implements Doku's B2B Access Token API
func (d *doku) GetToken(ctx context.Context) (*pg.TokenResponse, error) {
	// Validate private key is configured
	if d.config.PrivateKey == "" {
		return nil, fmt.Errorf("private key is required for GetToken API")
	}

	// Parse private key
	privateKey, err := d.parsePrivateKey(d.config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	baseURL := d.getBaseURL()
	fullURL := baseURL + tokenUri

	// Generate ISO8601 timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Build string to sign for asymmetric signature
	// Format: client_ID + "|" + timestamp
	stringToSign := d.config.ClientKey + "|" + timestamp

	// Generate SHA256withRSA signature
	signature, err := d.generateAsymmetricSignature(privateKey, stringToSign)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}

	// Build request body
	requestBody := map[string]string{
		"grantType": "client_credentials",
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for token API
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerXClientKey, d.config.ClientKey)
	req.Header.Set(headerXTimestamp, timestamp)
	req.Header.Set(headerXSignature, signature)

	// Execute request
	resp, err := d.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(responseBody))
	}

	// Parse token response
	var tokenResp tokenResponse
	if err := json.Unmarshal(responseBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &pg.TokenResponse{
		AccessToken: tokenResp.Token,
		TokenType:   "Bearer",
	}, nil
}

// parsePrivateKey parses a PEM-encoded RSA private key
func (d *doku) parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	// Try to parse as PKCS#1 or PKCS#8
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try PKCS#1 first
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Try PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaKey, nil
}

// generateAsymmetricSignature generates SHA256withRSA signature
func (d *doku) generateAsymmetricSignature(privateKey *rsa.PrivateKey, stringToSign string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(stringToSign))
	hashed := hasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// generateDigest generates SHA-256 digest of request body
// Digest format: base64(sha256(body))
func (d *doku) generateDigest(body []byte) string {
	hasher := sha256.New()
	hasher.Write(body)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// generateSignature generates Doku HMAC-SHA256 signature
// Signature component format:
// Client-Id:{clientId}
// Request-Id:{requestId}
// Request-Timestamp:{timestamp}
// Request-Target:{target}
// Digest:{digest}
func (d *doku) generateSignature(digest, timestamp, requestID, targetPath string) string {
	// Prepare signature component
	var component strings.Builder
	component.WriteString("Client-Id:" + d.config.ClientKey + "\n")
	component.WriteString("Request-Id:" + requestID + "\n")
	component.WriteString("Request-Timestamp:" + timestamp + "\n")
	component.WriteString("Request-Target:" + targetPath + "\n")
	component.WriteString("Digest:" + digest)

	// Calculate HMAC-SHA256
	h := hmac.New(sha256.New, []byte(d.config.ServerKey))
	h.Write([]byte(component.String()))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Prepend with algorithm info
	return "HMACSHA256=" + signature
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("req-%d-%s", time.Now().UnixMilli(), randomString(16))
}

// randomString generates a random string of specified length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

// validateChargeParams validates charge parameters
func (d *doku) validateChargeParams(params *pg.ChargeParams) error {
	// Validate order ID
	if err := utils.ValidateOrderID(params.OrderID); err != nil {
		return err
	}

	// Validate amount
	if err := utils.MinAmount(params.Amount, params.PaymentType); err != nil {
		return err
	}

	// Validate customer
	if params.Customer.ID == "" {
		return pg.NewRequiredFieldError("Customer.ID")
	}
	if err := utils.ValidateEmail(params.Customer.Email, "Customer.Email"); err != nil {
		return err
	}

	// Validate payment type
	if err := utils.ValidatePaymentType(params.PaymentType); err != nil {
		return err
	}

	return nil
}

// getTimeout returns the timeout duration
func getTimeout(cfg *pg.ProviderConfig) time.Duration {
	if cfg.Timeout > 0 {
		return time.Duration(cfg.Timeout) * time.Second
	}
	return 30 * time.Second
}

// getBaseURL returns the base URL based on environment
func (d *doku) getBaseURL() string {
	if d.config.Environment == "production" {
		return productionURL
	}
	return sandboxURL
}
