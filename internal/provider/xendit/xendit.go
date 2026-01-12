package xendit

import (
	"bytes"
	"context"
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
	// API endpoints
	invoiceUri     = "/v2/invoices"
	vaUri          = "/callback_virtual_accounts"
	ewalletUri     = "/ewallets"
	invoiceStatusUri = "/v2/invoices/%s"

	// header names
	headerAuthorization = "Authorization"
)

func init() {
	// Register this provider with the pg package
	pg.RegisterProvider(ProviderName, New)
}

type xendit struct {
	config   *pg.ProviderConfig
	mapper   *Mapper
	httpCli  *http.Client
}

// New creates a new Xendit provider
func New(cfg *pg.ProviderConfig) (pg.Provider, error) {
	if cfg.ServerKey == "" {
		return nil, pg.ErrMissingCredentials
	}

	return &xendit{
		config:  cfg,
		mapper:  &Mapper{},
		httpCli: &http.Client{Timeout: getTimeout(cfg)},
	}, nil
}

// Name returns the provider name
func (x *xendit) Name() string {
	return ProviderName
}

// CreateCharge creates a new payment transaction
func (x *xendit) CreateCharge(ctx context.Context, params pg.ChargeParams) (*pg.ChargeResponse, error) {
	// Validate parameters
	if err := x.validateChargeParams(&params); err != nil {
		return nil, err
	}

	var responseBody []byte
	var err error

	// Route to appropriate payment method
	if params.PaymentType.IsEWallet() || params.PaymentType == pg.PaymentTypeQRIS {
		ewalletReq := x.mapper.mapToEWalletRequest(params)
		responseBody, err = x.createEWalletCharge(ctx, ewalletReq)
		if err != nil {
			return nil, err
		}

		var resp EWalletResponse
		if err := json.Unmarshal(responseBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		return x.mapper.mapToChargeResponseFromEWallet(&resp, params.PaymentType), nil
	} else if params.PaymentType.IsVirtualAccount() {
		vaReq := x.mapper.mapToVARequest(params)
		responseBody, err = x.createVA(ctx, vaReq)
		if err != nil {
			return nil, err
		}

		var resp VAResponse
		if err := json.Unmarshal(responseBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		return x.mapper.mapToChargeResponseFromVA(&resp, params.PaymentType), nil
	} else {
		// Use Invoice API as fallback
		invoiceReq := x.mapper.mapToInvoiceRequest(params)
		responseBody, err = x.createInvoice(ctx, invoiceReq)
		if err != nil {
			return nil, err
		}

		var resp InvoiceResponse
		if err := json.Unmarshal(responseBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		return x.mapper.mapToChargeResponse(&resp, params.PaymentType), nil
	}
}

// createInvoice creates an invoice
func (x *xendit) createInvoice(ctx context.Context, params *CreateInvoiceRequest) ([]byte, error) {
	baseURL := x.getBaseURL()
	fullURL := baseURL + invoiceUri

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerAuthorization, utils.SetBasicAuthorization(x.config.ServerKey, ""))

	resp, err := x.httpCli.Do(req)
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

// createVA creates a virtual account
func (x *xendit) createVA(ctx context.Context, params *CreateVAResquest) ([]byte, error) {
	baseURL := x.getBaseURL()
	fullURL := baseURL + vaUri

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerAuthorization, utils.SetBasicAuthorization(x.config.ServerKey, ""))

	resp, err := x.httpCli.Do(req)
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

// createEWalletCharge creates an e-wallet charge
func (x *xendit) createEWalletCharge(ctx context.Context, params *CreateEWalletRequest) ([]byte, error) {
	baseURL := x.getBaseURL()
	fullURL := baseURL + ewalletUri

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerAuthorization, utils.SetBasicAuthorization(x.config.ServerKey, ""))

	resp, err := x.httpCli.Do(req)
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
func (x *xendit) GetStatus(ctx context.Context, orderID string) (*pg.PaymentStatus, error) {
	baseURL := x.getBaseURL()
	fullURL := fmt.Sprintf(baseURL+invoiceStatusUri, orderID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerAuthorization, utils.SetBasicAuthorization(x.config.ServerKey, ""))

	resp, err := x.httpCli.Do(req)
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

	var xenditResponse InvoiceResponse
	if err := json.Unmarshal(responseBody, &xenditResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return x.mapper.mapToPaymentStatus(orderID, &xenditResponse), nil
}

// Cancel cancels a transaction
func (x *xendit) Cancel(ctx context.Context, orderID string) error {
	baseURL := x.getBaseURL()
	fullURL := fmt.Sprintf(baseURL+invoiceStatusUri, orderID)

	// Xendit uses expiration to cancel (set expiration to past)
	expirePayload := map[string]interface{}{
		"expires_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
	}

	bodyBytes, err := json.Marshal(expirePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerAuthorization, utils.SetBasicAuthorization(x.config.ServerKey, ""))

	resp, err := x.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cancel failed: status=%d, body=%s", resp.StatusCode, string(responseBody))
	}

	return nil
}

// VerifyWebhook verifies webhook signature
func (x *xendit) VerifyWebhook(r *http.Request) bool {
	// Xendit uses X-Callback-Token header for webhook verification
	if x.config.ClientKey == "" {
		return false
	}

	verifier := utils.NewCallbackTokenVerifier(x.config.ClientKey)
	return verifier.Verify(r)
}

// GetToken retrieves an access token for the provider
// Xendit uses Basic Auth with API Key, so this returns an error
func (x *xendit) GetToken(ctx context.Context) (*pg.TokenResponse, error) {
	return nil, fmt.Errorf("GetToken API is not supported by %s. Xendit uses Basic Auth with API Key for authentication", ProviderName)
}

// ParseWebhook parses webhook payload
func (x *xendit) ParseWebhook(r *http.Request) (*pg.WebhookEvent, error) {
	if err := r.ParseForm(); err != nil {
		return nil, pg.ErrInvalidPayload
	}

	// Try to parse as JSON first
	var webhookData map[string]interface{}
	if r.Body != nil && r.ContentLength > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, pg.ErrInvalidPayload
		}
		if err := json.Unmarshal(body, &webhookData); err == nil {
			return x.parseWebhookJSON(webhookData)
		}
	}

	// Fall back to form parsing
	orderID := r.FormValue("external_id")
	status := r.FormValue("status")

	var amount int64
	if amountStr := r.FormValue("amount"); amountStr != "" {
		amount, _ = strconv.ParseInt(amountStr, 10, 64)
	}

	// Parse timestamp
	var timestamp time.Time
	if paidAt := r.FormValue("paid_at"); paidAt != "" {
		timestamp, _ = time.Parse(time.RFC3339, paidAt)
	}

	// Map status
	var mappedStatus pg.Status
	switch PaymentStatus(status) {
	case StatusPaid:
		mappedStatus = pg.StatusSuccess
	case StatusPending:
		mappedStatus = pg.StatusPending
	case StatusFailed:
		mappedStatus = pg.StatusFailed
	default:
		mappedStatus = pg.StatusPending
	}

	return &pg.WebhookEvent{
		OrderID:       orderID,
		TransactionID: orderID,
		Status:        mappedStatus,
		Amount:        amount,
		EventType:     x.mapEventType(PaymentStatus(status)),
		Timestamp:     timestamp,
		Raw:           getRawForm(r),
	}, nil
}

// parseWebhookJSON parses JSON webhook payload
func (x *xendit) parseWebhookJSON(data map[string]interface{}) (*pg.WebhookEvent, error) {
	orderID, _ := data["external_id"].(string)
	status, _ := data["status"].(string)

	var amount int64
	if amountFloat, ok := data["amount"].(float64); ok {
		amount = int64(amountFloat)
	}

	// Parse timestamp
	var timestamp time.Time
	if paidAt, ok := data["paid_at"].(string); ok {
		timestamp, _ = time.Parse(time.RFC3339, paidAt)
	}

	// Map status
	var mappedStatus pg.Status
	switch PaymentStatus(status) {
	case StatusPaid:
		mappedStatus = pg.StatusSuccess
	case StatusPending:
		mappedStatus = pg.StatusPending
	case StatusFailed:
		mappedStatus = pg.StatusFailed
	default:
		mappedStatus = pg.StatusPending
	}

	return &pg.WebhookEvent{
		OrderID:       orderID,
		TransactionID: orderID,
		Status:        mappedStatus,
		Amount:        amount,
		EventType:     x.mapEventType(PaymentStatus(status)),
		Timestamp:     timestamp,
		Raw:           data,
	}, nil
}

// getRawForm converts form data to raw map
func getRawForm(r *http.Request) map[string]interface{} {
	raw := make(map[string]interface{})
	r.ParseForm()
	for key, values := range r.Form {
		if len(values) == 1 {
			raw[key] = values[0]
		} else {
			raw[key] = values
		}
	}
	return raw
}

// mapEventType maps Xendit status to event type
func (x *xendit) mapEventType(status PaymentStatus) string {
	switch status {
	case StatusPaid:
		return pg.EventPaymentCompleted
	case StatusFailed:
		return pg.EventPaymentFailed
	case StatusPending:
		return pg.EventPaymentPending
	default:
		return pg.EventPaymentPending
	}
}

// validateChargeParams validates charge parameters
func (x *xendit) validateChargeParams(params *pg.ChargeParams) error {
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

// getBaseURL returns the base URL
func (x *xendit) getBaseURL() string {
	return sandboxURL // Xendit uses same URL for both environments
}
