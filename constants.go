package pg

// EnvironmentType represents the environment for payment gateway
type EnvironmentType string

const (
	// Production environment type for go live to real customer
	Production EnvironmentType = "production"

	// SandBox is staging (for development)
	SandBox EnvironmentType = "sandbox"
)

// PaymentType represents the type of payment method
type PaymentType string

// E-Wallet Payment Types
const (
	PaymentTypeGoPay     PaymentType = "GOPAY"
	PaymentTypeOVO       PaymentType = "OVO"
	PaymentTypeDANA      PaymentType = "DANA"
	PaymentTypeShopeePay PaymentType = "SHOPEEPAY"
	PaymentTypeLinkAja   PaymentType = "LINKAJA"
)

// Virtual Account Payment Types
const (
	PaymentTypeVABCA      PaymentType = "VA_BCA"
	PaymentTypeVABNI      PaymentType = "VA_BNI"
	PaymentTypeVABRI      PaymentType = "VA_BRI"
	PaymentTypeVAMandiri  PaymentType = "VA_MANDIRI"
	PaymentTypeVAPermata  PaymentType = "VA_PERMATA"
	PaymentTypeVACIMB     PaymentType = "VA_CIMB"
)

// QRIS Payment Type
const (
	PaymentTypeQRIS PaymentType = "QRIS"
)

// Credit Card Payment Type
const (
	PaymentTypeCC PaymentType = "CREDIT_CARD"
)

// Retail Payment Types
const (
	PaymentTypeAlfamart   PaymentType = "ALFAMART"
	PaymentTypeIndomaret  PaymentType = "INDOMARET"
)

// Status represents the status of a transaction
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusSuccess    Status = "SUCCESS"
	StatusFailed     Status = "FAILED"
	StatusCancelled  Status = "CANCELLED"
	StatusExpired    Status = "EXPIRED"
)

// Event Type for webhooks
const (
	EventPaymentCompleted = "payment.completed"
	EventPaymentFailed    = "payment.failed"
	EventPaymentPending   = "payment.pending"
	EventPaymentExpired   = "payment.expired"
	EventPaymentCancelled = "payment.cancelled"
)

// Provider names
const (
	ProviderMidtrans = "midtrans"
	ProviderXendit   = "xendit"
	ProviderDoku     = "doku"
)

// Provider names
const (
	ProviderDuitku  = "duitku"
	ProviderFaspay  = "faspay"
	ProviderEspay   = "espay"
)

// Default minimum amounts (in Rupiah)
const (
	MinAmountEWallet = 10000
	MinAmountVA      = 10000
	MinAmountQRIS    = 1000
	MinAmountCC      = 10000
	MinAmountRetail  = 10000
)

// Default expiry times
const (
	DefaultExpiryMinutes = 24 * 60 // 24 hours
)

// IsEWallet checks if the payment type is an e-wallet
func (p PaymentType) IsEWallet() bool {
	switch p {
	case PaymentTypeGoPay, PaymentTypeOVO, PaymentTypeDANA, PaymentTypeShopeePay, PaymentTypeLinkAja:
		return true
	}
	return false
}

// IsVirtualAccount checks if the payment type is a virtual account
func (p PaymentType) IsVirtualAccount() bool {
	switch p {
	case PaymentTypeVABCA, PaymentTypeVABNI, PaymentTypeVABRI, PaymentTypeVAMandiri,
		PaymentTypeVAPermata, PaymentTypeVACIMB:
		return true
	}
	return false
}

// IsQRIS checks if the payment type is QRIS
func (p PaymentType) IsQRIS() bool {
	return p == PaymentTypeQRIS
}

// IsCreditCard checks if the payment type is credit card
func (p PaymentType) IsCreditCard() bool {
	return p == PaymentTypeCC
}

// IsRetail checks if the payment type is retail outlet
func (p PaymentType) IsRetail() bool {
	switch p {
	case PaymentTypeAlfamart, PaymentTypeIndomaret:
		return true
	}
	return false
}

// String returns the string representation of the payment type
func (p PaymentType) String() string {
	return string(p)
}

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

// IsFinal returns true if the status is final (cannot be changed)
func (s Status) IsFinal() bool {
	switch s {
	case StatusSuccess, StatusFailed, StatusCancelled, StatusExpired:
		return true
	}
	return false
}

// String returns the string representation of the environment
func (e EnvironmentType) String() string {
	return string(e)
}

// IsProduction returns true if the environment is production
func (e EnvironmentType) IsProduction() bool {
	return e == Production
}

// Boolean variables for backward compatibility (must be variables to take address)
var (
	// True is a boolean true variable
	True = true
	// False is a boolean false variable
	False = false
)
