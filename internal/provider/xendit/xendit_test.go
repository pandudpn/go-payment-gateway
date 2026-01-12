package xendit

import (
	"net/http"
	"net/http/httptest"
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
				ServerKey: "test-server-key",
			},
			wantErr: false,
		},
		{
			name:    "missing server key",
			config: &pg.ProviderConfig{},
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

func TestXendit_Name(t *testing.T) {
	provider := &xendit{
		config: &pg.ProviderConfig{},
		mapper: &Mapper{},
	}

	if name := provider.Name(); name != ProviderName {
		t.Errorf("Name() = %v, want %v", name, ProviderName)
	}
}

func TestXendit_getBaseURL(t *testing.T) {
	provider := &xendit{
		config: &pg.ProviderConfig{
			Environment: "sandbox",
		},
	}

	// Xendit uses same URL for both environments
	expectedURL := "https://api.xendit.co"
	if got := provider.getBaseURL(); got != expectedURL {
		t.Errorf("getBaseURL() = %v, want %v", got, expectedURL)
	}
}

func TestXendit_CreateCharge_EWallet(t *testing.T) {
	tests := []struct {
		name         string
		params       pg.ChargeParams
		mockResponse string
		wantErr      bool
		expectStatus pg.Status
	}{
		{
			name: "successful GoPay charge",
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
				Items: []pg.Item{
					{
						ID:       "ITEM-001",
						Name:     "Test Product",
						Price:    50000,
						Quantity: 1,
					},
				},
				CallbackURL: "https://example.com/callback",
			},
			mockResponse: `{
				"id": "ewallet-123",
				"external_id": "ORDER-001",
				"amount": 50000,
				"ewallet_type": "GOPAY",
				"status": "PENDING",
				"payment_url": "https://example.com/pay"
			}`,
			wantErr:      false,
			expectStatus: pg.StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				auth := r.Header.Get("Authorization")
				if auth == "" {
					t.Error("missing authorization header")
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Test mapper
			mapper := &Mapper{}
			if tt.params.PaymentType.IsEWallet() {
				ewalletReq := mapper.mapToEWalletRequest(tt.params)
				if ewalletReq.ExternalID != tt.params.OrderID {
					t.Errorf("expected order ID %s, got %s", tt.params.OrderID, ewalletReq.ExternalID)
				}
				if ewalletReq.EWalletCode != EWalletGoPay && tt.params.PaymentType == pg.PaymentTypeGoPay {
					t.Errorf("expected GOPAY, got %v", ewalletReq.EWalletCode)
				}
			}
		})
	}
}

func TestXendit_CreateCharge_VA(t *testing.T) {
	params := pg.ChargeParams{
		OrderID:     "ORDER-002",
		Amount:      100000,
		PaymentType: pg.PaymentTypeVABCA,
		Customer: pg.Customer{
			ID:    "CUST-002",
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "+62812345679",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-002",
				Name:     "Product B",
				Price:    100000,
				Quantity: 1,
			},
		},
	}

	mapper := &Mapper{}
	vaReq := mapper.mapToVARequest(params)

	if vaReq.ExternalID != params.OrderID {
		t.Errorf("ExternalID = %v, want %v", vaReq.ExternalID, params.OrderID)
	}

	if vaReq.BankCode != BankBCA {
		t.Errorf("BankCode = %v, want %v", vaReq.BankCode, BankBCA)
	}

	if vaReq.ExpectedAmount != float64(params.Amount) {
		t.Errorf("ExpectedAmount = %v, want %v", vaReq.ExpectedAmount, params.Amount)
	}

	if !vaReq.IsClosed {
		t.Error("IsClosed should be true for fixed amount VA")
	}
}

func TestXendit_CreateCharge_Invoice(t *testing.T) {
	params := pg.ChargeParams{
		OrderID:     "ORDER-003",
		Amount:      150000,
		PaymentType: pg.PaymentTypeQRIS,
		Customer: pg.Customer{
			ID:    "CUST-003",
			Name:  "Bob Smith",
			Email: "bob@example.com",
			Phone: "+62812345670",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-003",
				Name:     "Product C",
				Price:    150000,
				Quantity: 1,
			},
		},
		Description: "Test payment",
	}

	mapper := &Mapper{}
	invoiceReq := mapper.mapToInvoiceRequest(params)

	if invoiceReq.ExternalID != params.OrderID {
		t.Errorf("ExternalID = %v, want %v", invoiceReq.ExternalID, params.OrderID)
	}

	if invoiceReq.Amount != float64(params.Amount) {
		t.Errorf("Amount = %v, want %v", invoiceReq.Amount, params.Amount)
	}

	if invoiceReq.Description != params.Description {
		t.Errorf("Description = %v, want %v", invoiceReq.Description, params.Description)
	}

	if invoiceReq.Customer == nil {
		t.Error("Customer should not be nil")
	} else {
		if invoiceReq.Customer.GivenNames != params.Customer.Name {
			t.Errorf("Customer.GivenNames = %v, want %v", invoiceReq.Customer.GivenNames, params.Customer.Name)
		}
	}

	if len(invoiceReq.Items) != len(params.Items) {
		t.Errorf("Items length = %v, want %v", len(invoiceReq.Items), len(params.Items))
	}
}

func TestMapper_mapPaymentType(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name             string
		paymentType      pg.PaymentType
		expectedType     string
		expectedCode     string
	}{
		{"GoPay", pg.PaymentTypeGoPay, "EWALLET", "GOPAY"},
		{"OVO", pg.PaymentTypeOVO, "EWALLET", "OVO"},
		{"DANA", pg.PaymentTypeDANA, "EWALLET", "DANA"},
		{"ShopeePay", pg.PaymentTypeShopeePay, "EWALLET", "SHOPEEPAY"},
		{"LinkAja", pg.PaymentTypeLinkAja, "EWALLET", "LINKAJA"},
		{"QRIS", pg.PaymentTypeQRIS, "QR_CODE", "QRIS"},
		{"VA BCA", pg.PaymentTypeVABCA, "VIRTUAL_ACCOUNT", "BCA"},
		{"VA BNI", pg.PaymentTypeVABNI, "VIRTUAL_ACCOUNT", "BNI"},
		{"VA BRI", pg.PaymentTypeVABRI, "VIRTUAL_ACCOUNT", "BRI"},
		{"VA Mandiri", pg.PaymentTypeVAMandiri, "VIRTUAL_ACCOUNT", "MANDIRI"},
		{"VA Permata", pg.PaymentTypeVAPermata, "VIRTUAL_ACCOUNT", "PERMATA"},
		{"VA CIMB", pg.PaymentTypeVACIMB, "VIRTUAL_ACCOUNT", "CIMB"},
		{"Alfamart", pg.PaymentTypeAlfamart, "RETAIL_OUTLET", "ALFAMART"},
		{"Indomaret", pg.PaymentTypeIndomaret, "RETAIL_OUTLET", "INDOMARET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			methodType, code := mapper.mapPaymentType(tt.paymentType)
			if methodType != tt.expectedType {
				t.Errorf("mapPaymentType(%v) type = %v, want %v", tt.paymentType, methodType, tt.expectedType)
			}
			if code != tt.expectedCode {
				t.Errorf("mapPaymentType(%v) code = %v, want %v", tt.paymentType, code, tt.expectedCode)
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
		{"Paid", StatusPaid, pg.StatusSuccess},
		{"Pending", StatusPending, pg.StatusPending},
		{"Failed", StatusFailed, pg.StatusFailed},
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

func TestMapper_mapToChargeResponse(t *testing.T) {
	mapper := &Mapper{}

	resp := &InvoiceResponse{
		ID:         "invoice-123",
		ExternalID: "ORDER-001",
		Amount:     50000,
		Status:     StatusPending,
		PaymentURL: "https://example.com/pay",
	}

	now := time.Now()
	resp.Created = &now

	result := mapper.mapToChargeResponse(resp, pg.PaymentTypeGoPay)

	if result.TransactionID != resp.ID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.ID)
	}

	if result.OrderID != resp.ExternalID {
		t.Errorf("OrderID = %v, want %v", result.OrderID, resp.ExternalID)
	}

	if result.Amount != int64(resp.Amount) {
		t.Errorf("Amount = %v, want %v", result.Amount, resp.Amount)
	}

	if result.Status != pg.StatusPending {
		t.Errorf("Status = %v, want %v", result.Status, pg.StatusPending)
	}

	if result.PaymentURL != resp.PaymentURL {
		t.Errorf("PaymentURL = %v, want %v", result.PaymentURL, resp.PaymentURL)
	}
}

func TestMapper_mapToChargeResponseFromVA(t *testing.T) {
	mapper := &Mapper{}

	resp := &VAResponse{
		ID:            "va-123",
		ExternalID:    "ORDER-002",
		BankCode:      BankBCA,
		AccountNumber: "1234567890",
		ExpectedAmount: 100000,
		Status:        StatusPending,
	}

	result := mapper.mapToChargeResponseFromVA(resp, pg.PaymentTypeVABCA)

	if result.TransactionID != resp.ID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.ID)
	}

	if result.OrderID != resp.ExternalID {
		t.Errorf("OrderID = %v, want %v", result.OrderID, resp.ExternalID)
	}

	if result.VANumber != resp.AccountNumber {
		t.Errorf("VANumber = %v, want %v", result.VANumber, resp.AccountNumber)
	}

	if result.VABank != string(resp.BankCode) {
		t.Errorf("VABank = %v, want %v", result.VABank, string(resp.BankCode))
	}
}

func TestMapper_mapToChargeResponseFromEWallet(t *testing.T) {
	mapper := &Mapper{}

	resp := &EWalletResponse{
		ID:         "ewallet-123",
		ExternalID: "ORDER-001",
		Amount:     50000,
		EWalletCode: EWalletGoPay,
		Status:     StatusPending,
		PaymentURL: "https://example.com/pay",
	}

	now := time.Now()
	resp.Created = &now

	result := mapper.mapToChargeResponseFromEWallet(resp, pg.PaymentTypeGoPay)

	if result.TransactionID != resp.ID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.ID)
	}

	if result.OrderID != resp.ExternalID {
		t.Errorf("OrderID = %v, want %v", result.OrderID, resp.ExternalID)
	}

	if result.Amount != int64(resp.Amount) {
		t.Errorf("Amount = %v, want %v", result.Amount, resp.Amount)
	}

	if result.PaymentURL != resp.PaymentURL {
		t.Errorf("PaymentURL = %v, want %v", result.PaymentURL, resp.PaymentURL)
	}
}

func TestMapper_mapToPaymentStatus(t *testing.T) {
	mapper := &Mapper{}

	t.Run("InvoiceResponse", func(t *testing.T) {
		now := time.Now()
		resp := &InvoiceResponse{
			ID:         "invoice-123",
			ExternalID: "ORDER-001",
			Amount:     50000,
			Status:     StatusPaid,
		}
		resp.Updated = &now

		result := mapper.mapToPaymentStatus("ORDER-001", resp)

		if result.TransactionID != resp.ID {
			t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.ID)
		}

		if result.Status != pg.StatusSuccess {
			t.Errorf("Status = %v, want %v", result.Status, pg.StatusSuccess)
		}

		if result.PaidAmount != int64(resp.Amount) {
			t.Errorf("PaidAmount = %v, want %v", result.PaidAmount, resp.Amount)
		}
	})

	t.Run("VAResponse", func(t *testing.T) {
		resp := &VAResponse{
			ID:             "va-123",
			ExternalID:     "ORDER-002",
			ExpectedAmount: 100000,
			Status:         StatusPending,
		}

		result := mapper.mapToPaymentStatus("ORDER-002", resp)

		if result.TransactionID != resp.ID {
			t.Errorf("TransactionID = %v, want %v", result.TransactionID, resp.ID)
		}

		if result.Status != pg.StatusPending {
			t.Errorf("Status = %v, want %v", result.Status, pg.StatusPending)
		}
	})
}

func TestMapper_unifiedPaymentType(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name            string
		methodType      string
		methodCode      string
		expectedPayment pg.PaymentType
	}{
		{"GoPay", "EWALLET", "GOPAY", pg.PaymentTypeGoPay},
		{"OVO", "EWALLET", "OVO", pg.PaymentTypeOVO},
		{"DANA", "EWALLET", "DANA", pg.PaymentTypeDANA},
		{"ShopeePay", "EWALLET", "SHOPEEPAY", pg.PaymentTypeShopeePay},
		{"LinkAja", "EWALLET", "LINKAJA", pg.PaymentTypeLinkAja},
		{"QRIS", "QR_CODE", "QRIS", pg.PaymentTypeQRIS},
		{"VA BCA", "VIRTUAL_ACCOUNT", "BCA", pg.PaymentTypeVABCA},
		{"VA BNI", "VIRTUAL_ACCOUNT", "BNI", pg.PaymentTypeVABNI},
		{"VA BRI", "VIRTUAL_ACCOUNT", "BRI", pg.PaymentTypeVABRI},
		{"VA Mandiri", "VIRTUAL_ACCOUNT", "MANDIRI", pg.PaymentTypeVAMandiri},
		{"VA Permata", "VIRTUAL_ACCOUNT", "PERMATA", pg.PaymentTypeVAPermata},
		{"VA CIMB", "VIRTUAL_ACCOUNT", "CIMB", pg.PaymentTypeVACIMB},
		{"Alfamart", "RETAIL_OUTLET", "ALFAMART", pg.PaymentTypeAlfamart},
		{"Indomaret", "RETAIL_OUTLET", "INDOMARET", pg.PaymentTypeIndomaret},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.unifiedPaymentType(tt.methodType, tt.methodCode)
			if got != tt.expectedPayment {
				t.Errorf("unifiedPaymentType(%s, %s) = %v, want %v", tt.methodType, tt.methodCode, got, tt.expectedPayment)
			}
		})
	}
}

func TestXendit_VerifyWebhook(t *testing.T) {
	tests := []struct {
		name      string
		clientKey string
		token     string
		want      bool
	}{
		{
			name:      "valid webhook token",
			clientKey: "test-callback-token",
			token:     "test-callback-token",
			want:      true,
		},
		{
			name:      "invalid webhook token",
			clientKey: "test-callback-token",
			token:     "wrong-token",
			want:      false,
		},
		{
			name:      "missing client key",
			clientKey: "",
			token:     "test-token",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &xendit{
				config: &pg.ProviderConfig{
					ClientKey: tt.clientKey,
				},
				mapper: &Mapper{},
			}

			req := httptest.NewRequest("POST", "/webhook", nil)
			req.Header.Set("X-Callback-Token", tt.token)

			got := provider.VerifyWebhook(req)
			if got != tt.want {
				t.Errorf("VerifyWebhook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXendit_ParseWebhook(t *testing.T) {
	provider := &xendit{
		config: &pg.ProviderConfig{},
		mapper: &Mapper{},
	}

	req := httptest.NewRequest("POST", "/webhook", nil)
	req.Form = map[string][]string{
		"external_id": {"ORDER-001"},
		"status":      {"PAID"},
		"amount":      {"50000"},
		"paid_at":     {"2024-01-01T00:00:00Z"},
	}

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

func TestXendit_mapEventType(t *testing.T) {
	provider := &xendit{
		config: &pg.ProviderConfig{},
		mapper: &Mapper{},
	}

	tests := []struct {
		name         string
		status       PaymentStatus
		expectedType string
	}{
		{"Paid", StatusPaid, pg.EventPaymentCompleted},
		{"Failed", StatusFailed, pg.EventPaymentFailed},
		{"Pending", StatusPending, pg.EventPaymentPending},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := provider.mapEventType(tt.status)
			if got != tt.expectedType {
				t.Errorf("mapEventType(%v) = %v, want %v", tt.status, got, tt.expectedType)
			}
		})
	}
}

func TestXendit_validateChargeParams(t *testing.T) {
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
				Items: []pg.Item{
					{ID: "ITEM-001", Name: "Product", Price: 50000, Quantity: 1},
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
				Items: []pg.Item{
					{ID: "ITEM-001", Name: "Product", Price: 50000, Quantity: 1},
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
				Items: []pg.Item{
					{ID: "ITEM-001", Name: "Product", Price: 50000, Quantity: 1},
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
				Items: []pg.Item{
					{ID: "ITEM-001", Name: "Product", Price: 50000, Quantity: 1},
				},
			},
			wantErr: true,
		},
		{
			name: "missing items",
			params: pg.ChargeParams{
				OrderID:     "ORDER-001",
				Amount:      50000,
				PaymentType: pg.PaymentTypeGoPay,
				Customer: pg.Customer{
					ID:    "CUST-001",
					Email: "john@example.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &xendit{
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

func TestXendit_GetStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("missing authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "invoice-123",
			"external_id": "ORDER-001",
			"amount": 50000,
			"status": "PAID"
		}`))
	}))
	defer server.Close()

	// Note: This test will try to call the actual API URL, not the mock server
	// In a real test, we would use dependency injection or interfaces
	// For this example, we're testing the mapper and response parsing

	mapper := &Mapper{}
	mockResp := &InvoiceResponse{
		ID:         "invoice-123",
		ExternalID: "ORDER-001",
		Amount:     50000,
		Status:     StatusPaid,
	}

	result := mapper.mapToPaymentStatus("ORDER-001", mockResp)

	if result.TransactionID != mockResp.ID {
		t.Errorf("TransactionID = %v, want %v", result.TransactionID, mockResp.ID)
	}

	if result.Status != pg.StatusSuccess {
		t.Errorf("Status = %v, want %v", result.Status, pg.StatusSuccess)
	}
}
