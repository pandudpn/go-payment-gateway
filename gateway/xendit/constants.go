package xendit

// Currency used for the transaction
type Currency string

const (
	// IDR is Indonesia Rupiah
	IDR Currency = "IDR"

	// PHP is Philippine Peso
	PHP Currency = "PHP"
)

// ItemType type of item or product in basket
type ItemType string

const (
	// Product item type
	Product ItemType = "PRODUCT"

	// Service item type
	Service ItemType = "SERVICE"
)

// ChannelCode specifies which e-wallet be used to process the transaction
type ChannelCode string

const (
	// ChannelCodeDANA is Payment Type or Channel Code for e-wallet DANA Indonesia from Xendit API
	// DANA supporting for payment:
	//  - One Time Payment
	//  - Linked Account (Tokenized)
	ChannelCodeDANA ChannelCode = "ID_DANA"

	// ChannelCodeOVO is Payment Type or Channel Code for e-wallet OVO from Xendit API
	// OVO supporting for payment:
	//  - One Time Payment
	//  - Linked Account (Tokenized)
	ChannelCodeOVO ChannelCode = "ID_OVO"

	// ChannelCodeShopeePay is Payment Type or Channel Code for e-wallet ShopeePay from Xendit API
	// ShopeePay supporting for payment:
	//  - One Time Payment
	//  - Linked Account (Tokenized)
	ChannelCodeShopeePay ChannelCode = "ID_SHOPEEPAY"

	// ChannelCodeLinkAja is Payment Type or Channel Code for e-wallet LinkAja from Xendit API
	// LinkAja only support One Time Payment
	ChannelCodeLinkAja ChannelCode = "ID_LINKAJA"
)

// CheckoutMethod determines the payment flow used to process the transaction
type CheckoutMethod string

const (
	// OneTimePayment is used for single guest checkouts
	OneTimePayment CheckoutMethod = "ONE_TIME_PAYMENT"

	// TokenizedPayment can be used for recurring payment or linked customer account
	TokenizedPayment CheckoutMethod = "TOKENIZED_PAYMENT"
)

// RedeemPoint is enum value when customer want to Use their Points in Tokenization CheckoutMethod
// used only for ChannelCode: EWalletOVO and EWalletShopeePay
type RedeemPoint string

const (
	// RedeemNone no points will be used in the transactions
	RedeemNone RedeemPoint = "REDEEM_NONE"

	// RedeemAll points will be used to offset payment amount before
	// cash balance is used
	// rules:
	// - for e-wallet ovo		: All points will be used
	// - for e-wallet shopee	: Only 50% of transaction amount (rounded down) can pay using shopeePay coins
	RedeemAll RedeemPoint = "REDEEM_ALL"
)

// ChargeStatus status of charge request
type ChargeStatus string

const (
	// Succeeded payment transaction for specified charge_id is successfully
	Succeeded ChargeStatus = "SUCCEEDED"

	// Pending payment transaction for specified charge_id is awaiting payment attempt by end user
	Pending ChargeStatus = "PENDING"

	// Failed payment transaction for specified charge_id has failed
	Failed ChargeStatus = "FAILED"

	// Voided payment transaction for specified charge_id has been voided
	Voided ChargeStatus = "VOIDED"

	// Refunded payment transaction for specified charge_id has been either partially or fully refunded
	Refunded ChargeStatus = "REFUNDED"
)
