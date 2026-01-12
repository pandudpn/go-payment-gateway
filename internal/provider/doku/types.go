package doku

import "time"

// PaymentStatus represents Doku payment status
type PaymentStatus string

const (
	// StatusPending means payment is pending
	StatusPending PaymentStatus = "PENDING"
	// StatusSuccess means payment is successful
	StatusSuccess PaymentStatus = "SUCCESS"
	// StatusFailed means payment failed
	StatusFailed PaymentStatus = "FAILED"
	// StatusCancelled means payment is cancelled
	StatusCancelled PaymentStatus = "CANCELLED"
)

// PaymentType represents Doku payment types
type PaymentType string

const (
	// PaymentTypeVirtualAccount for VA payments
	PaymentTypeVirtualAccount PaymentType = "VIRTUAL_ACCOUNT"
	// PaymentTypeEWallet for e-wallet payments
	PaymentTypeEWallet PaymentType = "EWALLET"
	// PaymentTypeQRCode for QR code payments
	PaymentTypeQRCode PaymentType = "QR_CODE"
	// PaymentTypePaylater for paylater payments
	PaymentTypePaylater PaymentType = "PAYLATER"
)

// VAComponent represents VA components for Doku
type VAComponent struct {
	Name          string `json:"name"`
	VaType        string `json:"va_type"`
	Amount        string `json:"amount"`
	CreatedDate   string `json:"created_date,omitempty"`
	ExpiredDate   string `json:"expired_date,omitempty"`
	VANumber      string `json:"virtual_account_number,omitempty"`
}

// EWalletComponent represents e-wallet components
type EWalletComponent struct {
	Name        string `json:"name"`
	EWalletType string `json:"e_wallet_type"`
	Amount      string `json:"amount"`
	Phone       string `json:"phone_number,omitempty"`
}

// QRCodeComponent represents QR code components
type QRCodeComponent struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
	QRType string `json:"qr_type"`
}

// PaylaterComponent represents paylater components
type PaylaterComponent struct {
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Tenure   int    `json:"tenure,omitempty"`
}

// GeneratePaymentRequest for Doku payment generation
type GeneratePaymentRequest struct {
	OrderAmount     int64           `json:"order_amount"`
	TransactionID   string          `json:"transaction_id"`
	TransactionDate time.Time       `json:"transaction_date"`
	Customer        *Customer       `json:"customer,omitempty"`
	PaymentType     PaymentType     `json:"payment_type"`
	PaymentDetail   *PaymentDetail  `json:"payment_detail,omitempty"`
	Locale          string          `json:"locale,omitempty"`
}

// Customer represents customer details
type Customer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	ID      string `json:"id,omitempty"`
}

// PaymentDetail represents payment details
type PaymentDetail struct {
	PaymentMethod   *PaymentMethod   `json:"payment_method,omitempty"`
	VirtualAccount  *VAComponent     `json:"virtual_account,omitempty"`
	EWallet         *EWalletComponent `json:"e_wallet,omitempty"`
	QRCode          *QRCodeComponent  `json:"qr_code,omitempty"`
	Paylater        *PaylaterComponent `json:"paylater,omitempty"`
}

// PaymentMethod represents payment method info
type PaymentMethod struct {
	Reus   bool `json:"reus,omitempty"`
	PayMethodDetails *PayMethodDetails `json:"pay_method_details,omitempty"`
}

// PayMethodDetails for additional payment method details
type PayMethodDetails struct {
	PayMethodCode string `json:"pay_method_code,omitempty"`
}

// GeneratePaymentResponse from Doku
type GeneratePaymentResponse struct {
	ResponseCode    string          `json:"response_code"`
	ResponseMessage string          `json:"response_message"`
	TransactionID   string          `json:"transaction_id,omitempty"`
	OrderAmount     int64           `json:"order_amount,omitempty"`
	PaymentURL      string          `json:"payment_url,omitempty"`
	VANumber        string          `json:"virtual_account_number,omitempty"`
	VABank          string          `json:"va_bank,omitempty"`
	QRString        string          `json:"qr_string,omitempty"`
}

// TransactionStatusRequest for checking transaction status
type TransactionStatusRequest struct {
	TransactionID string `json:"transaction_id"`
}

// TransactionStatusResponse from Doku
type TransactionStatusResponse struct {
	ResponseCode      string         `json:"response_code"`
	ResponseMessage   string         `json:"response_message"`
	TransactionID     string         `json:"transaction_id,omitempty"`
	OrderAmount       int64          `json:"order_amount,omitempty"`
	TransactionStatus PaymentStatus  `json:"transaction_status,omitempty"`
	PaymentType       PaymentType    `json:"payment_type,omitempty"`
	PaymentDate       *time.Time     `json:"payment_date,omitempty"`
}

// tokenResponse from Doku Get Token API
type tokenResponse struct {
	Token string `json:"token"`
}
