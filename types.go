package pg

import (
	"time"
)

// Customer represents customer information for payment
type Customer struct {
	// ID is the unique identifier for the customer
	ID string `json:"id"`

	// Name is the customer's full name
	Name string `json:"name"`

	// Email is the customer's email address
	Email string `json:"email"`

	// Phone is the customer's phone number
	Phone string `json:"phone"`
}

// Item represents a transaction item
type Item struct {
	// ID is the unique identifier for the item
	ID string `json:"id"`

	// Name is the item name
	Name string `json:"name"`

	// Price is the item price in smallest currency unit (e.g., Rupiah)
	Price int64 `json:"price"`

	// Quantity is the number of items
	Quantity int64 `json:"quantity"`

	// Category is the item category (optional)
	Category string `json:"category,omitempty"`

	// URL is the URL to the item page (optional)
	URL string `json:"url,omitempty"`
}

// ChargeParams represents the parameters for creating a payment charge
type ChargeParams struct {
	// OrderID is the unique identifier for the order (required)
	OrderID string `json:"order_id"`

	// Amount is the transaction amount in smallest currency unit (required)
	Amount int64 `json:"amount"`

	// PaymentType is the type of payment method (required)
	PaymentType PaymentType `json:"payment_type"`

	// Customer is the customer information (required)
	Customer Customer `json:"customer"`

	// Items is the list of transaction items
	Items []Item `json:"items,omitempty"`

	// Description is the transaction description
	Description string `json:"description,omitempty"`

	// ExpiryTime is when the payment will expire
	ExpiryTime time.Time `json:"expiry_time,omitempty"`

	// CallbackURL is the URL for payment status notifications
	CallbackURL string `json:"callback_url,omitempty"`

	// ReturnURL is the URL to redirect after payment
	ReturnURL string `json:"return_url,omitempty"`

	// Custom contains provider-specific parameters that are not mapped to unified fields
	// This allows access to provider-specific features
	Custom map[string]interface{} `json:"-"`
}

// ChargeResponse represents the response from creating a payment charge
type ChargeResponse struct {
	// TransactionID is the unique identifier from the payment provider
	TransactionID string `json:"transaction_id"`

	// OrderID is the merchant's order ID
	OrderID string `json:"order_id"`

	// Amount is the transaction amount
	Amount int64 `json:"amount"`

	// Status is the current payment status
	Status Status `json:"status"`

	// PaymentURL is the URL where user can complete the payment
	PaymentURL string `json:"payment_url,omitempty"`

	// QRString is the QR code string for QRIS payments
	QRString string `json:"qr_string,omitempty"`

	// VANumber is the virtual account number for VA payments
	VANumber string `json:"va_number,omitempty"`

	// VABank is the bank name for VA payments
	VABank string `json:"va_bank,omitempty"`

	// ExpiryTime is when the payment will expire
	ExpiryTime time.Time `json:"expiry_time,omitempty"`

	// CreatedAt is when the transaction was created
	CreatedAt time.Time `json:"created_at,omitempty"`

	// UpdatedAt is when the transaction was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// Raw contains the raw response from the provider for access to provider-specific fields
	Raw map[string]interface{} `json:"-"`
}

// PaymentStatus represents the status of a payment transaction
type PaymentStatus struct {
	// TransactionID is the unique identifier from the payment provider
	TransactionID string `json:"transaction_id"`

	// OrderID is the merchant's order ID
	OrderID string `json:"order_id"`

	// Status is the current payment status
	Status Status `json:"status"`

	// Amount is the transaction amount
	Amount int64 `json:"amount"`

	// PaidAmount is the amount that has been paid
	PaidAmount int64 `json:"paid_amount"`

	// PaymentType is the type of payment method used
	PaymentType PaymentType `json:"payment_type,omitempty"`

	// PaidAt is when the payment was completed
	PaidAt *time.Time `json:"paid_at,omitempty"`

	// CancelledAt is when the transaction was cancelled
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`

	// ExpiredAt is when the transaction expired
	ExpiredAt *time.Time `json:"expired_at,omitempty"`

	// FailureReason is the reason for payment failure (if applicable)
	FailureReason string `json:"failure_reason,omitempty"`

	// Raw contains the raw response from the provider
	Raw map[string]interface{} `json:"-"`
}

// WebhookEvent represents a webhook notification from the payment provider
type WebhookEvent struct {
	// OrderID is the merchant's order ID
	OrderID string `json:"order_id"`

	// TransactionID is the unique identifier from the payment provider
	TransactionID string `json:"transaction_id"`

	// Status is the payment status from the webhook
	Status Status `json:"status"`

	// Amount is the transaction amount
	Amount int64 `json:"amount"`

	// PaymentType is the type of payment method
	PaymentType PaymentType `json:"payment_type,omitempty"`

	// EventType is the type of webhook event
	EventType string `json:"event_type"`

	// Timestamp is when the webhook was sent
	Timestamp time.Time `json:"timestamp"`

	// FraudStatus is the fraud status (if applicable)
	FraudStatus string `json:"fraud_status,omitempty"`

	// Raw contains the raw webhook payload from the provider
	Raw map[string]interface{} `json:"-"`
}

// ProviderError represents an error from the payment provider
type ProviderError struct {
	// Code is the provider-specific error code
	Code string `json:"code"`

	// Message is the error message
	Message string `json:"message"`

	// Provider is the name of the payment provider
	Provider string `json:"provider"`

	// Raw contains the raw error response from the provider
	Raw map[string]interface{} `json:"-"`

	// Err is the underlying error
	Err error `json:"-"`
}

// Error returns the error message
func (e *ProviderError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "provider error"
}

// Unwrap returns the underlying error
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new ProviderError
func NewProviderError(provider, code, message string, raw map[string]interface{}) *ProviderError {
	return &ProviderError{
		Code:     code,
		Message:  message,
		Provider: provider,
		Raw:      raw,
	}
}

// CreditCardParams represents credit card specific parameters
type CreditCardParams struct {
	// CardNumber is the credit card number (tokenized)
	CardNumber string `json:"card_number"`

	// CardCvv is the card CVV (for 3DS)
	CardCvv string `json:"card_cvv,omitempty"`

	// CardExpMonth is the card expiration month (MM)
	CardExpMonth string `json:"card_exp_month,omitempty"`

	// CardExpYear is the card expiration year (YYYY)
	CardExpYear string `json:"card_exp_year,omitempty"`

	// SaveCard indicates if the card should be saved for future use
	SaveCard bool `json:"save_card,omitempty"`

	// Secure indicates if 3D Secure should be used
	Secure bool `json:"secure"`

	// InstallmentTerm is the installment term in months (0 for full payment)
	InstallmentTerm int `json:"installment_term,omitempty"`

	// Bank is the acquiring bank for installment
	Bank string `json:"bank,omitempty"`
}

// VirtualAccountParams represents virtual account specific parameters
type VirtualAccountParams struct {
	// BankCode is the bank code for the VA
	BankCode string `json:"bank_code"`

	// VANumber is a specific VA number (optional, will be generated if not provided)
	VANumber string `json:"va_number,omitempty"`

	// CustomerName is the customer name for VA display
	CustomerName string `json:"customer_name,omitempty"`
}

// EWALLETParams represents e-wallet specific parameters
type EWALLETParams struct {
	// RedirectToApp indicates if user should be redirected to e-wallet app
	RedirectToApp bool `json:"redirect_to_app,omitempty"`

	// AccountLinkID is the account link ID for recurring payments
	AccountLinkID string `json:"account_link_id,omitempty"`
}

// QRISParams represents QRIS specific parameters
type QRISParams struct {
	// QRString is the QR code string (generated by provider)
	QRString string `json:"qr_string,omitempty"`

	// Amount is a fixed amount QRIS (empty for dynamic QRIS)
	Amount int64 `json:"amount,omitempty"`
}

// TokenResponse represents the response from Get Token API
type TokenResponse struct {
	// AccessToken is the OAuth access token
	AccessToken string `json:"access_token"`

	// TokenType is the type of token (e.g., "Bearer")
	TokenType string `json:"token_type"`

	// ExpiresIn is the token expiration time in seconds
	ExpiresIn int64 `json:"expires_in"`

	// RefreshToken is the refresh token (if provided)
	RefreshToken string `json:"refresh_token,omitempty"`

	// Scope is the token scope
	Scope string `json:"scope,omitempty"`
}
