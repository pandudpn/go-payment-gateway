package doku

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	pg "github.com/pandudpn/go-payment-gateway"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *pg.ProviderConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &pg.ProviderConfig{
				ServerKey:  "test-server-key",
				ClientKey:  "test-client-key",
				Environment: "sandbox",
			},
			wantErr: false,
		},
		{
			name: "missing server key",
			config: &pg.ProviderConfig{
				ClientKey: "test-client-key",
			},
			wantErr: true,
		},
		{
			name: "missing client key",
			config: &pg.ProviderConfig{
				ServerKey: "test-server-key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("expected provider but got nil")
				}
			}
		})
	}
}

func TestDoku_Name(t *testing.T) {
	provider := &doku{
		config: &pg.ProviderConfig{},
		mapper: &Mapper{},
	}

	if name := provider.Name(); name != ProviderName {
		t.Errorf("Name() = %v, want %v", name, ProviderName)
	}
}

func TestDoku_getBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		wantURL  string
	}{
		{
			name:    "sandbox",
			env:     "sandbox",
			wantURL: "https://api-sandbox.doku.com",
		},
		{
			name:    "production",
			env:     "production",
			wantURL: "https://api.doku.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &doku{
				config: &pg.ProviderConfig{
					Environment: tt.env,
				},
			}

			if got := provider.getBaseURL(); got != tt.wantURL {
				t.Errorf("getBaseURL() = %v, want %v", got, tt.wantURL)
			}
		})
	}
}

func TestDoku_generateSignature(t *testing.T) {
	provider := &doku{
		config: &pg.ProviderConfig{
			ServerKey: "test-server-key",
			ClientKey: "test-client-id",
		},
		mapper: &Mapper{},
	}

	body := []byte(`{"test": "data"}`)
	timestamp := "1234567890"
	path := "/payments/v2"

	signature := provider.generateSignature(body, timestamp, path)

	// Verify signature format: CLIENT_ID:SIGNATURE
	parts := strings.Split(signature, ":")
	if len(parts) != 2 {
		t.Errorf("signature format incorrect, expected 2 parts, got %d", len(parts))
	}

	if parts[0] != provider.config.ClientKey {
		t.Errorf("signature client ID = %v, want %v", parts[0], provider.config.ClientKey)
	}

	// Verify the signature value
	digest := sha512.Sum512(body)
	hexDigest := hex.EncodeToString(digest[:])
	target := provider.config.ClientKey + ":" + timestamp + ":" + hexDigest

	h := hmac.New(sha512.New, []byte(provider.config.ServerKey))
	h.Write([]byte(target))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if parts[1] != expectedSignature {
		t.Errorf("signature value = %v, want %v", parts[1], expectedSignature)
	}
}

func TestDoku_VerifyWebhook(t *testing.T) {
	tests := []struct {
		name      string
		serverKey string
		clientKey string
		body      string
		signature string
		timestamp string
		want      bool
	}{
		{
			name:      "valid webhook signature",
			serverKey: "test-server-key",
			clientKey: "test-client-id",
			body:      `{"transaction_id": "ORDER-001"}`,
			timestamp: "1234567890",
			want:      true,
		},
		{
			name:      "invalid webhook signature",
			serverKey: "test-server-key",
			clientKey: "test-client-id",
			body:      `{"transaction_id": "ORDER-001"}`,
			signature: "invalid-signature",
			timestamp: "1234567890",
			want:      false,
		},
		{
			name:      "missing signature",
			serverKey: "test-server-key",
			clientKey: "test-client-id",
			body:      `{"transaction_id": "ORDER-001"}`,
			timestamp: "1234567890",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &doku{
				config: &pg.ProviderConfig{
					ServerKey: tt.serverKey,
					ClientKey: tt.clientKey,
				},
				mapper: &Mapper{},
			}

			// Calculate signature for valid test case
			if tt.want && tt.signature == "" {
				digest := sha512.Sum512([]byte(tt.body))
				hexDigest := hex.EncodeToString(digest[:])
				target := tt.clientKey + ":" + tt.timestamp + ":" + hexDigest
				h := hmac.New(sha512.New, []byte(tt.serverKey))
				h.Write([]byte(target))
				sig := hex.EncodeToString(h.Sum(nil))
				tt.signature = tt.clientKey + ":" + sig
			}

			req := httptest.NewRequest("POST", "/webhook", strings.NewReader(tt.body))
			req.Header.Set("Signature", tt.signature)
			req.Header.Set("Request-Timestamp", tt.timestamp)

			got := provider.VerifyWebhook(req)
			if got != tt.want {
				t.Errorf("VerifyWebhook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoku_ParseWebhook(t *testing.T) {
	provider := &doku{
		config: &pg.ProviderConfig{},
		mapper: &Mapper{},
	}

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(`{
		"transaction_id": "ORDER-001",
		"transaction_status": "SUCCESS",
		"amount": 50000,
		"payment_date_time": "2024-01-01T00:00:00Z"
	}`))
	req.Header.Set("Content-Type", "application/json")

	event, err := provider.ParseWebhook(req)
	if err != nil {
		t.Fatalf("ParseWebhook() error = %v", err)
	}

	if event.OrderID != "ORDER-001" {
		t.Errorf("OrderID = %v, want ORDER-001", event.OrderID)
	}

	if event.Status != pg.StatusSuccess {
		t.Errorf("Status = %v, want %v", event.Status, pg.StatusSuccess)
	}

	if event.Amount != 50000 {
		t.Errorf("Amount = %v, want 50000", event.Amount)
	}

	if event.EventType != pg.EventPaymentCompleted {
		t.Errorf("EventType = %v, want %v", event.EventType, pg.EventPaymentCompleted)
	}
}

func TestMapper_mapPaymentType(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name            string
		paymentType     pg.PaymentType
		expectedPayment PaymentType
	}{
		{"GoPay", pg.PaymentTypeGoPay, PaymentTypeEWallet},
		{"OVO", pg.PaymentTypeOVO, PaymentTypeEWallet},
		{"DANA", pg.PaymentTypeDANA, PaymentTypeEWallet},
		{"ShopeePay", pg.PaymentTypeShopeePay, PaymentTypeEWallet},
		{"LinkAja", pg.PaymentTypeLinkAja, PaymentTypeEWallet},
		{"QRIS", pg.PaymentTypeQRIS, PaymentTypeQRCode},
		{"VA BCA", pg.PaymentTypeVABCA, PaymentTypeVirtualAccount},
		{"VA BNI", pg.PaymentTypeVABNI, PaymentTypeVirtualAccount},
		{"VA BRI", pg.PaymentTypeVABRI, PaymentTypeVirtualAccount},
		{"VA Mandiri", pg.PaymentTypeVAMandiri, PaymentTypeVirtualAccount},
		{"VA Permata", pg.PaymentTypeVAPermata, PaymentTypeVirtualAccount},
		{"VA CIMB", pg.PaymentTypeVACIMB, PaymentTypeVirtualAccount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.mapPaymentType(tt.paymentType)
			if got != tt.expectedPayment {
				t.Errorf("mapPaymentType(%v) = %v, want %v", tt.paymentType, got, tt.expectedPayment)
			}
		})
	}
}

func TestMapper_mapStatus(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name           string
		status         PaymentStatus
		expectedStatus pg.Status
	}{
		{"Success", StatusSuccess, pg.StatusSuccess},
		{"Pending", StatusPending, pg.StatusPending},
		{"Failed", StatusFailed, pg.StatusFailed},
		{"Cancelled", StatusCancelled, pg.StatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.mapStatus(tt.status)
			if got != tt.expectedStatus {
				t.Errorf("mapStatus(%v) = %v, want %v", tt.status, got, tt.expectedStatus)
			}
		})
	}
}

func TestMapper_mapToGenerateRequest(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name        string
		params      pg.ChargeParams
		expectedPMT PaymentType
	}{
		{
			name: "GoPay e-wallet",
			params: pg.ChargeParams{
				OrderID:     "ORDER-001",
				Amount:      50000,
				PaymentType: pg.PaymentTypeGoPay,
				Customer: pg.Customer{
					ID:    "CUST-001",
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "+62812345678",
				},
			},
			expectedPMT: PaymentTypeEWallet,
		},
		{
			name: "QRIS",
			params: pg.ChargeParams{
				OrderID:     "ORDER-002",
				Amount:      100000,
				PaymentType: pg.PaymentTypeQRIS,
				Customer: pg.Customer{
					ID:    "CUST-002",
					Name:  "Jane Doe",
					Email: "jane@example.com",
				},
			},
			expectedPMT: PaymentTypeQRCode,
		},
		{
			name: "VA BCA",
			params: pg.ChargeParams{
				OrderID:     "ORDER-003",
				Amount:      150000,
				PaymentType: pg.PaymentTypeVABCA,
				Customer: pg.Customer{
					ID:    "CUST-003",
					Name:  "Bob Smith",
					Email: "bob@example.com",
				},
			},
			expectedPMT: PaymentTypeVirtualAccount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.mapToGenerateRequest(tt.params)

			if result.TransactionID != tt.params.OrderID {
				t.Errorf("TransactionID = %v, want %v", result.TransactionID, tt.params.OrderID)
			}

			if result.OrderAmount != tt.params.Amount {
				t.Errorf("OrderAmount = %v, want %v", result.OrderAmount, tt.params.Amount)
			}

			if result.PaymentType != tt.expectedPMT {
				t.Errorf("PaymentType = %v, want %v", result.PaymentType, tt.expectedPMT)
			}

			if result.Customer == nil {
				t.Error("Customer should not be nil")
			} else {
				if result.Customer.Name != tt.params.Customer.Name {
					t.Errorf("Customer.Name = %v, want %v", result.Customer.Name, tt.params.Customer.Name)
				}
				if result.Customer.Email != tt.params.Customer.Email {
					t.Errorf("Customer.Email = %v, want %v", result.Customer.Email, tt.params.Customer.Email)
				}
			}
		})
	}
}

func TestMapper_mapToChargeResponse(t *testing.T) {
	mapper := &Mapper{}

	resp := &GeneratePaymentResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		TransactionID:   "txn-123",
		OrderAmount:     50000,
		PaymentURL:      "https://example.com/pay",
		VANumber:        "1234567890",
		VABank:          "BCA",
	}

	result := mapper.mapToChargeResponse(resp, pg.PaymentTypeVABCA)

	if result.TransactionID != resp.TransactionID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.TransactionID)
	}

	if result.OrderID != resp.TransactionID {
		t.Errorf("OrderID = %v, want %v", result.OrderID, resp.TransactionID)
	}

	if result.Amount != resp.OrderAmount {
		t.Errorf("Amount = %v, want %v", result.Amount, resp.OrderAmount)
	}

	if result.Status != pg.StatusPending {
		t.Errorf("Status = %v, want %v", result.Status, pg.StatusPending)
	}

	if result.PaymentURL != resp.PaymentURL {
		t.Errorf("PaymentURL = %v, want %v", result.PaymentURL, resp.PaymentURL)
	}

	if result.VANumber != resp.VANumber {
		t.Errorf("VANumber = %v, want %v", result.VANumber, resp.VANumber)
	}

	if result.VABank != resp.VABank {
		t.Errorf("VABank = %v, want %v", result.VABank, resp.VABank)
	}
}

func TestMapper_mapToPaymentStatus(t *testing.T) {
	mapper := &Mapper{}

	now := time.Now()
	resp := &TransactionStatusResponse{
		ResponseCode:      "00",
		ResponseMessage:   "Success",
		TransactionID:     "txn-123",
		OrderAmount:       50000,
		TransactionStatus: StatusSuccess,
		PaymentDate:       &now,
	}

	result := mapper.mapToPaymentStatus("ORDER-001", resp)

	if result.TransactionID != resp.TransactionID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.TransactionID)
	}

	if result.OrderID != "ORDER-001" {
		t.Errorf("OrderID = %v, want ORDER-001", result.OrderID)
	}

	if result.Status != pg.StatusSuccess {
		t.Errorf("Status = %v, want %v", result.Status, pg.StatusSuccess)
	}

	if result.Amount != resp.OrderAmount {
		t.Errorf("Amount = %v, want %v", result.Amount, resp.OrderAmount)
	}

	if result.PaidAmount != resp.OrderAmount {
		t.Errorf("PaidAmount = %v, want %v", result.PaidAmount, resp.OrderAmount)
	}

	if result.PaidAt == nil {
		t.Error("PaidAt should not be nil for successful payment")
	}
}

func TestMapper_mapEventType(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name         string
		status       PaymentStatus
		expectedType string
	}{
		{"Success", StatusSuccess, pg.EventPaymentCompleted},
		{"Failed", StatusFailed, pg.EventPaymentFailed},
		{"Cancelled", StatusCancelled, pg.EventPaymentCancelled},
		{"Pending", StatusPending, pg.EventPaymentPending},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.mapEventType(tt.status)
			if got != tt.expectedType {
				t.Errorf("mapEventType(%v) = %v, want %v", tt.status, got, tt.expectedType)
			}
		})
	}
}

func TestDoku_validateChargeParams(t *testing.T) {
	tests := []struct {
		name    string
		params  pg.ChargeParams
		wantErr bool
	}{
		{
			name: "valid params",
			params: pg.ChargeParams{
				OrderID:     "ORDER-001",
				Amount:      50000,
				PaymentType: pg.PaymentTypeGoPay,
				Customer: pg.Customer{
					ID:    "CUST-001",
					Email: "john@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "missing order ID",
			params: pg.ChargeParams{
				Amount:      50000,
				PaymentType: pg.PaymentTypeGoPay,
				Customer: pg.Customer{
					ID:    "CUST-001",
					Email: "john@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "missing customer ID",
			params: pg.ChargeParams{
				OrderID:     "ORDER-001",
				Amount:      50000,
				PaymentType: pg.PaymentTypeGoPay,
				Customer: pg.Customer{
					Email: "john@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			params: pg.ChargeParams{
				OrderID:     "ORDER-001",
				Amount:      50000,
				PaymentType: pg.PaymentTypeGoPay,
				Customer: pg.Customer{
					ID:    "CUST-001",
					Email: "invalid-email",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &doku{
				config: &pg.ProviderConfig{},
				mapper: &Mapper{},
			}

			err := provider.validateChargeParams(&tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateChargeParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDoku_CreateCharge(t *testing.T) {
	tests := []struct {
		name         string
		params       pg.ChargeParams
		mockResponse string
		wantErr      bool
	}{
		{
			name: "successful VA charge",
			params: pg.ChargeParams{
				OrderID:     "ORDER-001",
				Amount:      100000,
				PaymentType: pg.PaymentTypeVABCA,
				Customer: pg.Customer{
					ID:    "CUST-001",
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "+62812345678",
				},
			},
			mockResponse: `{
				"response_code": "00",
				"response_message": "Success",
				"transaction_id": "txn-123",
				"order_amount": 100000,
				"virtual_account_number": "1234567890",
				"va_bank": "BCA"
			}`,
			wantErr: false,
		},
		{
			name: "successful QRIS charge",
			params: pg.ChargeParams{
				OrderID:     "ORDER-002",
				Amount:      50000,
				PaymentType: pg.PaymentTypeQRIS,
				Customer: pg.Customer{
					ID:    "CUST-002",
					Name:  "Jane Doe",
					Email: "jane@example.com",
				},
			},
			mockResponse: `{
				"response_code": "00",
				"response_message": "Success",
				"transaction_id": "txn-456",
				"order_amount": 50000,
				"payment_url": "https://example.com/qr"
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				// Check required headers
				clientID := r.Header.Get("Client-Id")
				if clientID == "" {
					t.Error("missing Client-Id header")
				}

				timestamp := r.Header.Get("Request-Timestamp")
				if timestamp == "" {
					t.Error("missing Request-Timestamp header")
				}

				signature := r.Header.Get("Signature")
				if signature == "" {
					t.Error("missing Signature header")
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Test mapper
			mapper := &Mapper{}
			req := mapper.mapToGenerateRequest(tt.params)

			if req.TransactionID != tt.params.OrderID {
				t.Errorf("TransactionID = %v, want %v", req.TransactionID, tt.params.OrderID)
			}

			if req.OrderAmount != tt.params.Amount {
				t.Errorf("OrderAmount = %v, want %v", req.OrderAmount, tt.params.Amount)
			}
		})
	}
}

func TestDoku_GetStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"response_code": "00",
			"response_message": "Success",
			"transaction_id": "txn-123",
			"order_amount": 50000,
			"transaction_status": "SUCCESS"
		}`))
	}))
	defer server.Close()

	// Test mapper
	mapper := &Mapper{}
	mockResp := &TransactionStatusResponse{
		ResponseCode:      "00",
		TransactionID:     "txn-123",
		OrderAmount:       50000,
		TransactionStatus: StatusSuccess,
	}

	result := mapper.mapToPaymentStatus("ORDER-001", mockResp)

	if result.TransactionID != mockResp.TransactionID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, mockResp.TransactionID)
	}

	if result.Status != pg.StatusSuccess {
		t.Errorf("Status = %v, want %v", result.Status, pg.StatusSuccess)
	}
}

func TestDoku_Cancel(t *testing.T) {
	provider := &doku{
		config: &pg.ProviderConfig{},
		mapper: &Mapper{},
	}

	err := provider.Cancel(context.Background(), "ORDER-001")
	if err == nil {
		t.Error("expected error for cancel, got nil")
	}

	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("error message should mention not supported, got: %v", err)
	}
}
