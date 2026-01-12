package xendit

import "time"

// PaymentStatus represents Xendit payment status
type PaymentStatus string

const (
	// StatusPending means payment is pending
	StatusPending PaymentStatus = "PENDING"
	// StatusPaid means payment is paid
	StatusPaid PaymentStatus = "PAID"
	// StatusFailed means payment failed
	StatusFailed PaymentStatus = "FAILED"
)

// PaymentChannel represents Xendit payment channels
type PaymentChannel string

const (
	// ChannelVirtualAccount for VA payments
	ChannelVirtualAccount PaymentChannel = "VIRTUAL_ACCOUNT"
	// ChannelEWallet for e-wallet payments
	ChannelEWallet PaymentChannel = "EWALLET"
	// ChannelQRCode for QR code payments
	ChannelQRCode PaymentChannel = "QR_CODE"
	// ChannelRetail for retail payments
	ChannelRetail PaymentChannel = "RETAIL_OUTLET"
)

// BankCode represents Xendit bank codes for VA
type BankCode string

const (
	// BankBCA Xendit VA BCA
	BankBCA BankCode = "BCA"
	// BankBNI Xendit VA BNI
	BankBNI BankCode = "BNI"
	// BankBRI Xendit VA BRI
	BankBRI BankCode = "BRI"
	// BankMANDIRI Xendit VA Mandiri
	BankMANDIRI BankCode = "MANDIRI"
	// BankPERMATA Xendit VA Permata
	BankPERMATA BankCode = "PERMATA"
	// BankCIMB Xendit VA CIMB
	BankCIMB BankCode = "CIMB"
	// BankSahabatSampoerna Xendit VA Sampoerna
	BankSahabatSampoerna BankCode = "SAHABAT_SAMPURNA"
)

// EWalletCode represents Xendit e-wallet codes
type EWalletCode string

const (
	// EWalletGoPay Xendit GoPay
	EWalletGoPay EWalletCode = "GOPAY"
	// EWalletOVO Xendit OVO
	EWalletOVO EWalletCode = "OVO"
	// EWalletDANA Xendit DANA
	EWalletDANA EWalletCode = "DANA"
	// EWalletLinkAja Xendit LinkAja
	EWalletLinkAja EWalletCode = "LINKAJA"
	// EWalletShopeePay Xendit ShopeePay
	EWalletShopeePay EWalletCode = "SHOPEEPAY"
)

// QRCodeCode represents Xendit QR code types
type QRCodeCode string

const (
	// QRCodeDynamic Xendit QRIS Dynamic
	QRCodeDynamic QRCodeCode = "QRIS"
)

// RetailCode represents Xendit retail outlet codes
type RetailCode string

const (
	// RetailAlfamart Xendit Alfamart
	RetailAlfamart RetailCode = "ALFAMART"
	// RetailIndomaret Xendit Indomaret
	RetailIndomaret RetailCode = "INDOMARET"
)

// CreateInvoiceRequest for creating Xendit invoice
type CreateInvoiceRequest struct {
	ExternalID       string              `json:"external_id"`
	Amount           float64             `json:"amount"`
	InvoiceDuration   int64               `json:"invoice_duration,omitempty"`
	Description       string              `json:"description,omitempty"`
	PaymentMethod     []PaymentMethod     `json:"payment_methods,omitempty"`
	Currency          string              `json:"currency,omitempty"`
	ReminderTime      int64               `json:"reminder_time,omitempty"`
	Customer          *CustomerDetail     `json:"customer,omitempty"`
	CustomerDetails   *CustomerNotification `json:"customer_notification_preference,omitempty"`
	Items             []*Item             `json:"items,omitempty"`
	Fees              []*Fee              `json:"fees,omitempty"`
	ShouldSendEmail   bool                `json:"should_send_email,omitempty"`
	ForUserID         string              `json:"for_user_id,omitempty"`
	Platform          *Platform           `json:"platform,omitempty"`
}

// PaymentMethod represents payment method details
type PaymentMethod struct {
	Type    string      `json:"type,omitempty"`
	Reus    bool        `json:"reus,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// CustomerDetail for Xendit
type CustomerDetail struct {
	GivenNames string `json:"given_names,omitempty"`
	Email      string `json:"email,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	CustomerID string `json:"customer_id,omitempty"`
}

// CustomerNotification for notifications
type CustomerNotification struct {
	Mails   []*string `json:" mails,omitempty"`
	Locale  string    `json:"locale,omitempty"`
}

// Item represents item details
type Item struct {
	Name        string  `json:"name,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Quantity    int64   `json:"quantity,omitempty"`
	Category    string  `json:"category,omitempty"`
	URL         string  `json:"url,omitempty"`
}

// Fee represents fee details
type Fee struct {
	Type        string  `json:"type,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Description string  `json:"description,omitempty"`
}

// Platform details
type Platform struct {
	ReconciliationDetail *ReconciliationDetail `json:"reconciliation_detail,omitempty"`
}

// ReconciliationDetail for reconciliation
type ReconciliationDetail struct {
	Type string `json:"type,omitempty"`
}

// InvoiceResponse from Xendit
type InvoiceResponse struct {
	ID                  string          `json:"id"`
	ExternalID          string          `json:"external_id"`
	Amount              float64         `json:"amount"`
	Status              PaymentStatus   `json:"status"`
	Currency            string          `json:"currency,omitempty"`
	BusinessID          string          `json:"business_id,omitempty"`
	Created             *time.Time      `json:"created,omitempty"`
	Updated             *time.Time      `json:"updated,omitempty"`
	PaymentURL          string          `json:"payment_url,omitempty"`
	ExpirationDate      *time.Time      `json:"expiration_date,omitempty"`
	InvoiceURL          string          `json:"invoice_url,omitempty"`
	UserID              string          `json:"user_id,omitempty"`
	PaymentMethod       string          `json:"payment_method,omitempty"`
	PaymentChannel      string          `json:"payment_channel,omitempty"`
	PaymentDestination  string          `json:"payment_destination,omitempty"`
	PaymentDetails      *PaymentDetails  `json:"payment_details,omitempty"`
	Customer            *CustomerDetail `json:"customer,omitempty"`
	Items               []*Item         `json:"items,omitempty"`
	MerchantProfile     *MerchantProfile `json:"merchant_profile_profile,omitempty"`
	Metadata            map[string]string `json:"metadata,omitempty"`
}

// PaymentDetails for invoice response
type PaymentDetails struct {
	Source               string `json:"source,omitempty"`
	Destination          string `json:"destination,omitempty"`
	ReceiptNotification   string `json:"receipt_notification,omitempty"`
	BillerCode           string `json:"biller_code,omitempty"`
	BillerKey            string `json:"biller_key,omitempty"`
	PaymentID            string `json:"payment_id,omitempty"`
	Amount               float64 `json:"amount,omitempty"`
}

// MerchantProfile details
type MerchantProfile struct {
	BusinessName string `json:"business_name,omitempty"`
}

// CreateVAResquest for creating VA
type CreateVAResquest struct {
	ExternalID       string   `json:"external_id"`
	BankCode         BankCode `json:"bank_code"`
	Name             string   `json:"name"`
	ExpectedAmount   float64  `json:"expected_amount,omitempty"`
	IsClosed         bool     `json:"is_closed,omitempty"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
	Currency         string   `json:"currency,omitempty"`
	VANumber         string   `json:"virtual_account_number,omitempty"`
	SuggestedAmount  bool     `json:"suggested_amount,omitempty"`
	IsSingleUse      bool     `json:"is_single_use,omitempty"`
	Description      string   `json:"description,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// VAResponse from Xendit
type VAResponse struct {
	ID                string          `json:"id"`
	ExternalID        string          `json:"external_id"`
	BankCode          BankCode        `json:"bank_code"`
	MerchantCode      string          `json:"merchant_code,omitempty"`
	Name              string          `json:"name"`
	AccountNumber     string          `json:"account_number,omitempty"`
	VANumber          string          `json:"virtual_account_number,omitempty"`
	ExpectedAmount    float64         `json:"expected_amount,omitempty"`
	IsClosed          bool            `json:"is_closed,omitempty"`
	ExpirationDate    *time.Time      `json:"expiration_date,omitempty"`
	Currency          string          `json:"currency,omitempty"`
	Status            PaymentStatus   `json:"status,omitempty"`
	SuggestedAmount   bool            `json:"suggested_amount,omitempty"`
	IsSingleUse       bool            `json:"is_single_use,omitempty"`
	Description       string          `json:"description,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	Payment           *PaymentDetails  `json:"payment,omitempty"`
}

// CreateEWalletRequest for creating e-wallet payment
type CreateEWalletRequest struct {
	ExternalID       string                `json:"external_id"`
	Amount           float64               `json:"amount"`
	Phone            string                `json:"phone,omitempty"`
	CallbackURL      string                `json:"callback_url,omitempty"`
	RedirectURL      string                `json:"redirect_url,omitempty"`
	EWalletCode      EWalletCode           `json:"ewallet_type"`
	BusinessID       string                `json:"batch_id,omitempty"`
	Currency         string                `json:"currency,omitempty"`
	ChannelProperties *ChannelProperties    `json:"channel_properties,omitempty"`
	Metadata         map[string]string     `json:"metadata,omitempty"`
}

// ChannelProperties for e-wallet
type ChannelProperties struct {
	SuccessRedirectURL string `json:"success_redirect_url,omitempty"`
	FailureRedirectURL string `json:"failure_redirect_url,omitempty"`
	PendingRedirectURL string `json:"pending_redirect_url,omitempty"`
	 RedeemPoints      bool   `json:"redeem_points,omitempty"`
}

// EWalletResponse from Xendit
type EWalletResponse struct {
	ID                string              `json:"id,omitempty"`
	ExternalID        string              `json:"external_id,omitempty"`
	Amount            float64             `json:"amount,omitempty"`
	EWalletCode       EWalletCode         `json:"ewallet_type,omitempty"`
	Status            PaymentStatus       `json:"status,omitempty"`
	CallbackURL       string              `json:"callback_url,omitempty"`
	PaymentURL        string              `json:"payment_url,omitempty"`
	RedirectURL       string              `json:"redirect_url,omitempty"`
	Created           *time.Time          `json:"created,omitempty"`
	Updated           *time.Time          `json:"updated,omitempty"`
	Metadata          map[string]string   `json:"metadata,omitempty"`
}
