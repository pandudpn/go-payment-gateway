package midtrans

import (
	"context"
	"encoding/json"
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
	// API endpoints
	chargeUri    = "/v2/charge"
	statusUri    = "/v2/%s/status"
	cancelUri    = "/v2/%s/cancel"

	// header names
	headerAuthorization = "Authorization"
)

type midtrans struct {
	config   *pg.ProviderConfig
	mapper   *Mapper
	httpCli  *http.Client
}

// New creates a new Midtrans provider
func New(cfg *pg.ProviderConfig) (pg.Provider, error) {
	if cfg.ServerKey == "" {
		return nil, pg.ErrMissingCredentials
	}

	return &midtrans{
		config:  cfg,
		mapper:  &Mapper{},
		httpCli: &http.Client{Timeout: getTimeout(cfg)},
	}, nil
}

// Name returns the provider name
func (m *midtrans) Name() string {
	return ProviderName
}

// CreateCharge creates a new payment transaction
func (m *midtrans) CreateCharge(ctx context.Context, params pg.ChargeParams) (*pg.ChargeResponse, error) {
	// Validate parameters
	if err := m.validateChargeParams(&params); err != nil {
		return nil, err
	}

	// Map to Midtrans params
	var responseBody []byte
	var err error

	if params.PaymentType.IsEWallet() || params.PaymentType == pg.PaymentTypeQRIS {
		ewalletParams := m.mapper.mapToEWalletParams(params)
		responseBody, err = m.createChargeEWallet(ctx, ewalletParams)
	} else if params.PaymentType.IsVirtualAccount() {
		bankParams := m.mapper.mapToBankTransferParams(params)
		responseBody, err = m.createChargeBankTransfer(ctx, bankParams)
	} else {
		return nil, pg.NewFieldError("PaymentType", fmt.Sprintf("payment type %s is not yet supported", params.PaymentType))
	}

	if err != nil {
		return nil, err
	}

	// Parse response
	var midtransResponse ChargeResponse
	if err := json.Unmarshal(responseBody, &midtransResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for error in response
	if midtransResponse.StatusCode != "201" && midtransResponse.StatusCode != "200" {
		return nil, pg.WrapProviderError(ProviderName, midtransResponse.StatusCode, midtransResponse.StatusMessage, nil)
	}

	// Map to unified response
	return m.mapper.mapToChargeResponse(&midtransResponse), nil
}

// createChargeEWallet creates e-wallet charge
func (m *midtrans) createChargeEWallet(ctx context.Context, params *EWallet) ([]byte, error) {
	baseURL := m.getBaseURL()

	// Build request
	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+chargeUri, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", utils.SetBasicAuthorization(m.config.ServerKey, ""))

	// Execute request
	resp, err := m.httpCli.Do(req)
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

// createChargeBankTransfer creates bank transfer charge
func (m *midtrans) createChargeBankTransfer(ctx context.Context, params *BankTransferCreateParams) ([]byte, error) {
	baseURL := m.getBaseURL()

	// Build request
	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+chargeUri, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", utils.SetBasicAuthorization(m.config.ServerKey, ""))

	// Execute request
	resp, err := m.httpCli.Do(req)
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
func (m *midtrans) GetStatus(ctx context.Context, orderID string) (*pg.PaymentStatus, error) {
	baseURL := m.getBaseURL()
	fullURL := fmt.Sprintf(baseURL+statusUri, orderID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", utils.SetBasicAuthorization(m.config.ServerKey, ""))

	resp, err := m.httpCli.Do(req)
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

	var midtransResponse ChargeResponse
	if err := json.Unmarshal(responseBody, &midtransResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return m.mapper.mapToPaymentStatus(orderID, &midtransResponse), nil
}

// Cancel cancels a transaction
func (m *midtrans) Cancel(ctx context.Context, orderID string) error {
	baseURL := m.getBaseURL()
	fullURL := fmt.Sprintf(baseURL+cancelUri, orderID)

	// Build cancel request
	cancelPayload := map[string]string{
		"cancel_reason": "User requested cancellation",
	}

	bodyBytes, err := json.Marshal(cancelPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", utils.SetBasicAuthorization(m.config.ServerKey, ""))

	resp, err := m.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("cancel failed: status=%d, body=%s", resp.StatusCode, string(responseBody))
	}

	return nil
}

// VerifyWebhook verifies webhook signature
func (m *midtrans) VerifyWebhook(r *http.Request) bool {
	// Get order ID and status from request
	if err := r.ParseForm(); err != nil {
		return false
	}

	orderID := r.FormValue("order_id")
	status := r.FormValue("transaction_status")

	if orderID == "" || status == "" {
		return false
	}

	// Create signature verifier
	verifier := utils.NewSHA512Verifier(m.config.ServerKey, "order_id", "transaction_status")

	return verifier.VerifyString(orderID, status, r.Header.Get("X-Signature"))
}

// ParseWebhook parses webhook payload
func (m *midtrans) ParseWebhook(r *http.Request) (*pg.WebhookEvent, error) {
	if err := r.ParseForm(); err != nil {
		return nil, pg.ErrInvalidPayload
	}

	orderID := r.FormValue("order_id")
	transactionStatus := r.FormValue("transaction_status")
	fraudStatus := r.FormValue("fraud_status")

	// Parse gross amount
	var amount int64
	if grossAmount := r.FormValue("gross_amount"); grossAmount != "" {
		amount, _ = strconv.ParseInt(grossAmount, 10, 64)
	}

	// Parse transaction time
	var timestamp time.Time
	if transactionTime := r.FormValue("transaction_time"); transactionTime != "" {
		timestamp, _ = time.Parse(time.RFC3339, transactionTime)
	}

	// Map status
	var mappedStatus pg.Status
	switch TransactionStatus(transactionStatus) {
	case Settlement, Capture:
		mappedStatus = pg.StatusSuccess
	case Pending:
		mappedStatus = pg.StatusPending
	case Deny, Failure:
		mappedStatus = pg.StatusFailed
	case Cancel:
		mappedStatus = pg.StatusCancelled
	case Expire:
		mappedStatus = pg.StatusExpired
	default:
		mappedStatus = pg.StatusPending
	}

	// Build raw response map
	raw := make(map[string]interface{})
	r.ParseForm()
	for key, values := range r.Form {
		if len(values) == 1 {
			raw[key] = values[0]
		} else {
			raw[key] = values
		}
	}

	return &pg.WebhookEvent{
		OrderID:       orderID,
		TransactionID: orderID, // Midtrans uses order_id as both
		Status:        mappedStatus,
		Amount:        amount,
		EventType:     m.mapper.mapEventType(TransactionStatus(transactionStatus)),
		Timestamp:     timestamp,
		FraudStatus:    fraudStatus,
		Raw:           raw,
	}, nil
}

// validateChargeParams validates charge parameters
func (m *midtrans) validateChargeParams(params *pg.ChargeParams) error {
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
	if err := utils.ValidatePhone(params.Customer.Phone, "Customer.Phone"); err != nil {
		return err
	}

	// Validate payment type
	if err := utils.ValidatePaymentType(params.PaymentType); err != nil {
		return err
	}

	// Validate items
	if len(params.Items) == 0 {
		return pg.NewRequiredFieldError("Items")
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
func (m *midtrans) getBaseURL() string {
	if m.config.SnapMode {
		if m.config.Environment == "production" {
			return snapProductionURL
		}
		return snapSandboxURL
	}

	if m.config.Environment == "production" {
		return productionURL
	}
	return sandboxURL
}
