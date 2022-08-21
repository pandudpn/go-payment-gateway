package midtrans

// TransactionStatus status of payment transaction
type TransactionStatus string

const (
	// Authorize means the payment card used for the transaction
	// must be captured to process the balance
	Authorize TransactionStatus = "authorize"

	// Capture means transaction is success and card balance is captured successfully
	Capture TransactionStatus = "capture"

	// Settlement means transaction is successfully settled.
	Settlement TransactionStatus = "settlement"

	// Pending means transaction still wait to payment
	Pending TransactionStatus = "pending"

	// Deny means payments are rejected by the payment provider
	// or Midtrans Fraud Detection System
	Deny TransactionStatus = "deny"

	// Cancel means the transaction is cancelled.
	Cancel TransactionStatus = "cancel"

	// Refund the transaction is marked to be refunded.
	Refund TransactionStatus = "refund"

	// PartialRefund the transaction is marked to be partially refunded
	PartialRefund TransactionStatus = "partial_refund"

	// Chargeback the transaction is marked to be charged back
	Chargeback TransactionStatus = "chargeback"

	// PartialChargeback the transaction is marked to be partially charged back
	PartialChargeback TransactionStatus = "partial_chargeback"

	// Expire the transaction is not available for processing. because the payment was delayed
	Expire TransactionStatus = "expire"

	// Failure means unexpected error occurred during transaction processing
	Failure TransactionStatus = "failure"
)

type PaymentType string

// list midtrans payment type
const (
	// PaymentTypeGopay is payment type gopay from Midtrans Core API
	// gopay supporting for payment:
	//  - One Time Payment
	//  - Linked Account (Tokenized)
	//  - Recurring
	PaymentTypeGopay PaymentType = "gopay"

	// PaymentTypeShopeePay is payment type ShopeePay from Midtrans Core API
	// shopeePay from Midtrans only support One Time Payment
	PaymentTypeShopeePay PaymentType = "shopeepay"

	// PaymentTypeBCA is payment type for Virtual Account with Bank BCA from Midtrans Core API
	// bca from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	PaymentTypeBCA PaymentType = "bank_transfer"

	// PaymentTypeBRI is payment type for Virtual Account with Bank BRI from Midtrans Core API
	// bri from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	PaymentTypeBRI PaymentType = "bank_transfer"

	// PaymentTypeBNI is payment type for Virtual Account with Bank BNI from Midtrans Core API
	// bni from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	PaymentTypeBNI PaymentType = "bank_transfer"

	// PaymentTypeMandiri is payment type for Virtual Account with Bank Mandiri from Midtrans Core API
	// mandiri from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	PaymentTypeMandiri PaymentType = "echannel"

	// PaymentTypePermata is payment type for Virtual Account with Bank Permata from Midtrans Core API
	// permata from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	PaymentTypePermata PaymentType = "bank_transfer"

	// PaymentTypeCard is payment type for Credit Card or Debit Card from Midtrans Core API
	PaymentTypeCard PaymentType = "credit_card"

	// PaymentTypeAkulaku is payment type for paylater Akulaku from Midtrans Core API
	PaymentTypeAkulaku PaymentType = "akulaku"

	// PaymentTypeKredivo is payment type for paylater Kredivo from Midtrans Core API
	PaymentTypeKredivo PaymentType = "kredivo"
)

// convert to String value or type
func (pt PaymentType) String() string {
	return string(pt)
}

type BankCode string

const (
	// BankBCA code for Bank BCA
	BankBCA BankCode = "bca"

	// BankBRI code for Bank BRI
	BankBRI BankCode = "bri"

	// BankBNI code for Bank BNI
	BankBNI BankCode = "bni"

	// BankPermata code for Bank Permata
	BankPermata BankCode = "permata"

	// BankMandiri code for Bank Mandiri
	BankMandiri BankCode = "mandiri"

	// BankMaybank code for Bank Maybank
	BankMaybank BankCode = "maybank"

	// BankCIMB code for Bank CIMB Niaga
	BankCIMB BankCode = "cimb"
)

// Unit of expiry duration
type ExpiryUnitDuration string

const (
	// Second payment will expired in second
	Second ExpiryUnitDuration = "second"

	// Minute payment will expired in minute
	Minute ExpiryUnitDuration = "minute"

	// Hour payment will expired in hour
	Hour ExpiryUnitDuration = "hour"

	// Day payment will expired in day
	Day ExpiryUnitDuration = "day"
)

type FraudStatus string

const (
	// FraudAccept means approved by FDS
	FraudAccept FraudStatus = "accept"

	// FraudChallenge means questioned by FDS
	FraudChallenge FraudStatus = "challenge"

	// FraudDeny denied by FDS. transaction automatically failed
	FraudDeny FraudStatus = "deny"
)

type CardType string

const (
	// Debit card
	Debit CardType = "debit"

	// Credit card
	Credit CardType = "credit"
)

type AccountStatus string

const (
	// AccountPending means account waiting to linked
	AccountPending AccountStatus = "PENDING"

	// AccountEnabled success to linked
	AccountEnabled AccountStatus = "ENABLED"

	// AccountExpired time to linked is expired
	AccountExpired AccountStatus = "EXPIRED"

	// AccountDisabled account already to unlinked
	AccountDisabled AccountStatus = "DISABLED"
)
