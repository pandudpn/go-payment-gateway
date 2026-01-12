package pg

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomer_JSONMarshal(t *testing.T) {
	customer := Customer{
		ID:    "CUST-001",
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+628123456789",
	}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := `{"id":"CUST-001","name":"John Doe","email":"john@example.com","phone":"+628123456789"}`
	if string(data) != expected {
		t.Errorf("Marshal() = %v, want %v", string(data), expected)
	}
}

func TestItem_JSONMarshal(t *testing.T) {
	item := Item{
		ID:       "ITEM-001",
		Name:     "Product A",
		Price:    50000,
		Quantity: 2,
		Category: "Electronics",
		URL:      "https://example.com/product",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Check that all fields are present
	str := string(data)
	if !contains(str, `"id":"ITEM-001"`) {
		t.Error("ID field not found in JSON")
	}
	if !contains(str, `"name":"Product A"`) {
		t.Error("Name field not found in JSON")
	}
	if !contains(str, `"price":50000`) {
		t.Error("Price field not found in JSON")
	}
	if !contains(str, `"quantity":2`) {
		t.Error("Quantity field not found in JSON")
	}
	if !contains(str, `"category":"Electronics"`) {
		t.Error("Category field not found in JSON")
	}
	if !contains(str, `"url":"https://example.com/product"`) {
		t.Error("URL field not found in JSON")
	}
}

func TestChargeParams_JSONMarshal(t *testing.T) {
	params := ChargeParams{
		OrderID:     "ORDER-001",
		Amount:      100000,
		PaymentType: PaymentTypeGoPay,
		Customer: Customer{
			ID:    "CUST-001",
			Name:  "John Doe",
			Email: "john@example.com",
		},
		Items: []Item{
			{
				ID:       "ITEM-001",
				Name:     "Product A",
				Price:    50000,
				Quantity: 2,
			},
		},
		Description: "Test payment",
		CallbackURL: "https://example.com/callback",
		ReturnURL:   "https://example.com/return",
		Custom: map[string]interface{}{
			"custom_field": "custom_value",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Check that Custom is not marshaled (it has json:"-" tag)
	str := string(data)
	if contains(str, "custom_field") {
		t.Error("Custom field should not be marshaled (has json:\"-\")")
	}
}

func TestChargeResponse_JSONMarshal(t *testing.T) {
	now := time.Now()
	resp := ChargeResponse{
		TransactionID: "txn-123",
		OrderID:       "ORDER-001",
		Amount:        100000,
		Status:        StatusSuccess,
		PaymentURL:    "https://example.com/pay",
		VANumber:      "1234567890",
		VABank:        "BCA",
		CreatedAt:     now,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Check that Raw is not marshaled
	str := string(data)
	if contains(str, `"raw":`) {
		t.Error("Raw field should not be marshaled (has json:\"-\")")
	}

	// Check important fields are present
	if !contains(str, `"transaction_id":"txn-123"`) {
		t.Error("transaction_id field not found")
	}
	if !contains(str, `"status":"SUCCESS"`) {
		t.Error("status field not found")
	}
}

func TestPaymentStatus_JSONMarshal(t *testing.T) {
	now := time.Now()
	status := PaymentStatus{
		TransactionID: "txn-123",
		OrderID:       "ORDER-001",
		Status:        StatusSuccess,
		Amount:        100000,
		PaidAmount:    100000,
		PaymentType:   PaymentTypeGoPay,
		PaidAt:        &now,
		FailureReason: "test failure",
	}

	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	str := string(data)
	if contains(str, `"raw":`) {
		t.Error("Raw field should not be marshaled")
	}

	if !contains(str, `"status":"SUCCESS"`) {
		t.Error("status field not found")
	}
}

func TestWebhookEvent_JSONMarshal(t *testing.T) {
	event := WebhookEvent{
		OrderID:       "ORDER-001",
		TransactionID: "txn-123",
		Status:        StatusSuccess,
		Amount:        50000,
		PaymentType:   PaymentTypeGoPay,
		EventType:     EventPaymentCompleted,
		Timestamp:     time.Now(),
		FraudStatus:   "accept",
		Raw:           map[string]interface{}{"raw": "data"},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	str := string(data)
	if contains(str, `"raw":`) {
		t.Error("Raw field should not be marshaled")
	}

	if !contains(str, `"event_type":"payment.completed"`) {
		t.Error("event_type field not found")
	}
}

func TestProviderError_JSONMarshal(t *testing.T) {
	providerErr := &ProviderError{
		Code:     "400",
		Message:  "Bad Request",
		Provider: "midtrans",
		Raw:      map[string]interface{}{"status": "error"},
	}

	data, err := json.Marshal(providerErr)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	str := string(data)
	if contains(str, `"raw":`) {
		t.Error("Raw field should not be marshaled")
	}
	if contains(str, `"err":`) {
		t.Error("Err field should not be marshaled")
	}
}

func TestCreditCardParams(t *testing.T) {
	params := CreditCardParams{
		CardNumber:     "4111111111111111",
		CardCvv:        "123",
		CardExpMonth:   "12",
		CardExpYear:    "2025",
		SaveCard:       true,
		Secure:         true,
		InstallmentTerm: 12,
		Bank:           "bca",
	}

	if params.CardNumber != "4111111111111111" {
		t.Error("CardNumber not set correctly")
	}

	if !params.SaveCard {
		t.Error("SaveCard should be true")
	}

	if !params.Secure {
		t.Error("Secure should be true")
	}

	if params.InstallmentTerm != 12 {
		t.Errorf("InstallmentTerm = %v, want 12", params.InstallmentTerm)
	}
}

func TestVirtualAccountParams(t *testing.T) {
	params := VirtualAccountParams{
		BankCode:     "BCA",
		VANumber:     "1234567890",
		CustomerName: "John Doe",
	}

	if params.BankCode != "BCA" {
		t.Error("BankCode not set correctly")
	}

	if params.VANumber != "1234567890" {
		t.Error("VANumber not set correctly")
	}

	if params.CustomerName != "John Doe" {
		t.Error("CustomerName not set correctly")
	}
}

func TestEWALLETParams(t *testing.T) {
	params := EWALLETParams{
		RedirectToApp:  true,
		AccountLinkID:  "link-123",
	}

	if !params.RedirectToApp {
		t.Error("RedirectToApp should be true")
	}

	if params.AccountLinkID != "link-123" {
		t.Error("AccountLinkID not set correctly")
	}
}

func TestQRISParams(t *testing.T) {
	params := QRISParams{
		QRString: "00020101021234567890",
		Amount:   50000,
	}

	if params.QRString != "00020101021234567890" {
		t.Error("QRString not set correctly")
	}

	if params.Amount != 50000 {
		t.Errorf("Amount = %v, want 50000", params.Amount)
	}
}

func TestChargeResponse_AllFields(t *testing.T) {
	now := time.Now()
	resp := ChargeResponse{
		TransactionID: "txn-123",
		OrderID:       "ORDER-001",
		Amount:        50000,
		Status:        StatusPending,
		PaymentURL:    "https://pay.example.com",
		QRString:      "qr-string-123",
		VANumber:      "888999",
		VABank:        "BCA",
		ExpiryTime:    now.Add(24 * time.Hour),
		CreatedAt:     now,
		UpdatedAt:     now,
		Raw:           map[string]interface{}{"provider_data": "value"},
	}

	if resp.TransactionID != "txn-123" {
		t.Errorf("TransactionID = %v, want txn-123", resp.TransactionID)
	}

	if resp.OrderID != "ORDER-001" {
		t.Errorf("OrderID = %v, want ORDER-001", resp.OrderID)
	}

	if resp.Amount != 50000 {
		t.Errorf("Amount = %v, want 50000", resp.Amount)
	}

	if resp.Status != StatusPending {
		t.Errorf("Status = %v, want PENDING", resp.Status)
	}

	if resp.QRString != "qr-string-123" {
		t.Errorf("QRString = %v, want qr-string-123", resp.QRString)
	}

	if resp.VANumber != "888999" {
		t.Errorf("VANumber = %v, want 888999", resp.VANumber)
	}

	if resp.VABank != "BCA" {
		t.Errorf("VABank = %v, want BCA", resp.VABank)
	}
}

func TestPaymentStatus_AllFields(t *testing.T) {
	now := time.Now()
	paidAt := now.Add(1 * time.Hour)
	status := PaymentStatus{
		TransactionID: "txn-456",
		OrderID:       "ORDER-002",
		Status:        StatusSuccess,
		Amount:        100000,
		PaidAmount:    100000,
		PaymentType:   PaymentTypeVABCA,
		PaidAt:        &paidAt,
		FailureReason: "none",
		Raw:           map[string]interface{}{"data": "value"},
	}

	if status.TransactionID != "txn-456" {
		t.Errorf("TransactionID = %v, want txn-456", status.TransactionID)
	}

	if status.Status != StatusSuccess {
		t.Errorf("Status = %v, want SUCCESS", status.Status)
	}

	if status.PaidAt == nil {
		t.Error("PaidAt should not be nil")
	}

	if status.PaidAmount != 100000 {
		t.Errorf("PaidAmount = %v, want 100000", status.PaidAmount)
	}
}

func TestWebhookEvent_AllFields(t *testing.T) {
	now := time.Now()
	event := WebhookEvent{
		OrderID:       "ORDER-003",
		TransactionID: "txn-789",
		Status:        StatusFailed,
		Amount:        75000,
		PaymentType:   PaymentTypeOVO,
		EventType:     EventPaymentFailed,
		Timestamp:     now,
		FraudStatus:   "reject",
		Raw:           map[string]interface{}{"webhook": "data"},
	}

	if event.EventType != EventPaymentFailed {
		t.Errorf("EventType = %v, want %v", event.EventType, EventPaymentFailed)
	}

	if event.FraudStatus != "reject" {
		t.Errorf("FraudStatus = %v, want reject", event.FraudStatus)
	}

	if event.Amount != 75000 {
		t.Errorf("Amount = %v, want 75000", event.Amount)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
