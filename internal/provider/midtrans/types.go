package midtrans

import "time"

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

	// PaymentTypeOVO is payment type OVO from Midtrans Core API
	PaymentTypeOVO PaymentType = "ovo"

	// PaymentTypeDANA is payment type DANA from Midtrans Core API
	PaymentTypeDANA PaymentType = "dana"

	// PaymentTypeLinkAja is payment type LinkAja from Midtrans Core API
	PaymentTypeLinkAja PaymentType = "linkaja"

	// PaymentTypeQRIS is payment type QRIS from Midtrans Core API
	PaymentTypeQRIS PaymentType = "qris"

	// PaymentTypeBCA is payment type for Virtual Account with Bank BCA from Midtrans Core API
	PaymentTypeBCA PaymentType = "bank_transfer"

	// PaymentTypeBRI is payment type for Virtual Account with Bank BRI from Midtrans Core API
	PaymentTypeBRI PaymentType = "bank_transfer"

	// PaymentTypeBNI is payment type for Virtual Account with Bank BNI from Midtrans Core API
	PaymentTypeBNI PaymentType = "bank_transfer"

	// PaymentTypeMandiri is payment type for Virtual Account with Bank Mandiri from Midtrans Core API
	PaymentTypeMandiri PaymentType = "echannel"

	// PaymentTypePermata is payment type for Virtual Account with Bank Permata from Midtrans Core API
	PaymentTypePermata PaymentType = "bank_transfer"

	// PaymentTypeCard is payment type for Credit Card or Debit Card from Midtrans Core API
	PaymentTypeCard PaymentType = "credit_card"
)

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

// TransactionDetail for customer
type TransactionDetail struct {
	// OrderID is Reference ID for midtrans
	OrderID string `json:"order_id"`

	// GrossAmount is total transaction to be paid by Customer
	GrossAmount int64 `json:"gross_amount"`
}

// CustomerDetail details of customer
type CustomerDetail struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// ItemDetail details items purchased by Customer
type ItemDetail struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
	Category string `json:"category,omitempty"`
}

// EWalletDetail details e-wallet payment
type EWalletDetail struct {
	CallbackURL string `json:"callback_url,omitempty"`
	AccountID   string `json:"account_id,omitempty"`
}

// EWallet define payload for e-wallet charge
type EWallet struct {
	PaymentType       PaymentType       `json:"payment_type"`
	TransactionDetails *TransactionDetail `json:"transaction_details"`
	ItemDetails       []*ItemDetail      `json:"item_details"`
	CustomerDetails   *CustomerDetail   `json:"customer_details,omitempty"`
	Gopay             *EWalletDetail    `json:"gopay,omitempty"`
	ShopeePay          *EWalletDetail    `json:"shopeepay,omitempty"`
	OVO                *EWalletDetail    `json:"ovo,omitempty"`
	DANA               *EWalletDetail    `json:"dana,omitempty"`
	LinkAja            *EWalletDetail    `json:"linkaja,omitempty"`
}

// BankTransfer charge details using bank transfer
type BankTransfer struct {
	Bank     BankCode `json:"bank"`
	VANumber string   `json:"va_number,omitempty"`
}

// EChannel charge details using Mandiri Bill Payment
type EChannel struct {
	BillInfo1 string `json:"bill_info1"`
	BillInfo2 string `json:"bill_info2"`
}

// BankTransferCreateParams parameters for VA charge
type BankTransferCreateParams struct {
	PaymentType       PaymentType       `json:"payment_type"`
	TransactionDetails *TransactionDetail `json:"transaction_details"`
	ItemDetails       []*ItemDetail      `json:"item_details"`
	CustomerDetails   *CustomerDetail   `json:"customer_details,omitempty"`
	BankTransfer      *BankTransfer      `json:"bank_transfer,omitempty"`
	EChannel          *EChannel          `json:"echannel,omitempty"`
}

// Action to make payments redirect
type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// ChargeResponse charge response from Midtrans
type ChargeResponse struct {
	StatusCode           string           `json:"status_code"`
	StatusMessage        string           `json:"status_message"`
	TransactionID        string           `json:"transaction_id"`
	OrderID              string           `json:"order_id"`
	GrossAmount          string           `json:"gross_amount"`
	PaymentType          PaymentType      `json:"payment_type"`
	TransactionTime      time.Time        `json:"transaction_time"`
	TransactionStatus    TransactionStatus `json:"transaction_status"`
	FraudStatus          string           `json:"fraud_status"`
	RedirectURL          string           `json:"redirect_url"`
	Actions              []*Action         `json:"actions"`
	BillKey              string           `json:"bill_key"`
	BillerCode           string           `json:"biller_code"`
	PermataVANumber      string           `json:"permata_va_number"`
	VANumbers            []*BankTransfer  `json:"va_numbers"`
	Bank                 BankCode         `json:"bank"`
}
