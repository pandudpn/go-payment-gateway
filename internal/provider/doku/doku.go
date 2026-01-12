package doku

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

	// header names
	headerSignature = "Signature"
	headerTimestamp = "Request-Timestamp"
	headerClientID  = "Client-Id"

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

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	signature := d.generateSignature(bodyBytes, timestamp, "/payments/v2")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerClientID, d.config.ClientKey)
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

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	signature := d.generateSignature(bodyBytes, timestamp, "/transactions/v2")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set(headerClientID, d.config.ClientKey)
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

	if signature == "" || timestamp == "" {
		return false
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// Restore body for subsequent reads
	r.Body = io.NopCloser(bytes.NewReader(body))

	// Calculate expected signature
	digest := sha512.Sum512(body)
	hexDigest := hex.EncodeToString(digest[:])

	target := fmt.Sprintf("%s:%s:%s", d.config.ClientKey, timestamp, hexDigest)
	h := hmac.New(sha512.New, []byte(d.config.ServerKey))
	h.Write([]byte(target))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Signature header format: CLIENT_ID:SIGNATURE
	// We need to extract the signature part and compare, or reconstruct full format
	expectedFullSignature := fmt.Sprintf("%s:%s", d.config.ClientKey, expectedSignature)

	return signature == expectedFullSignature
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

// generateSignature generates Doku signature
func (d *doku) generateSignature(body []byte, timestamp, path string) string {
	digest := sha512.Sum512(body)
	hexDigest := hex.EncodeToString(digest[:])

	// Doku signature format: CLIENT_ID:TIMESTAMP:REQUEST_BODY_DIGEST
	target := fmt.Sprintf("%s:%s:%s", d.config.ClientKey, timestamp, hexDigest)

	// HMAC-SHA512 with Server Key
	h := hmac.New(sha512.New, []byte(d.config.ServerKey))
	h.Write([]byte(target))
	signature := hex.EncodeToString(h.Sum(nil))

	// Final signature: CLIENT_ID:SIGNATURE
	return fmt.Sprintf("%s:%s", d.config.ClientKey, signature)
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
