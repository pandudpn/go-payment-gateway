package xendit

import (
	"time"
)

// EWalletChannelProperties specific information required for the e-wallet transaction to be initiated
type EWalletChannelProperties struct {
	// MobileNumber of customer in E.164 format (e.g: 628123456789)
	// this field required when customer choose OVO as a payments
	//
	// Default: ""
	MobileNumber string `json:"mobile_number,omitempty"`

	// SuccessRedirectURL URL where the end-customer is redirected
	// if the authorization is successful
	// this field required when customer choose OVO, DANA, LinkAja, ShopeePay as a Payments
	//
	// Default: ""
	SuccessRedirectURL string `json:"success_redirect_url,omitempty"`

	// FailureRedirectURL url where the end-customer is redirected
	// if the authorization is failed
	// this field is required when customer choose OVO as a Payments
	//
	// Default: ""
	FailureRedirectURL string `json:"failure_redirect_url,omitempty"`

	// RedeemPoints customer can use their point when CheckoutMethod is TokenizedPayment
	//
	// Default: ""
	RedeemPoints RedeemPoint `json:"redeem_points,omitempty"`
}

// EWalletItem describing the items purchased by customer
type EWalletItem struct {
	// ReferenceID merchant's identifier for specific product
	// this field required
	//
	// Default: ""
	ReferenceID string `json:"reference_id"`

	// Name of product
	// this field required
	//
	// Default: ""
	Name string `json:"name"`

	// Category of product
	// this field required
	//
	// Default: ""
	Category string `json:"category"`

	// Currency used for the transaction
	// this field required and possible values only: IDR or PHP
	//
	// Default: "IDR"
	Currency Currency `json:"currency"`

	// Price per unit in item currency
	// this field required
	//
	// Default: 0
	Price float64 `json:"price"`

	// Quantity number of units of this item
	// this field required
	//
	// Default: 0
	Quantity float64 `json:"quantity"`

	// Type of product
	// this field required and possible values only: PRODUCT or SERVICE
	//
	// Default: "PRODUCT"
	Type ItemType `json:"type"`

	// URL to page of item
	//
	// Default: ""
	URL string `json:"url,omitempty"`

	// Description of product
	//
	// Default: ""
	Description string `json:"description,omitempty"`

	// SubCategory for item
	//
	// Default: ""
	SubCategory string `json:"sub_category,omitempty"`

	// Metadata object of additional information the user may use for this item
	//
	// default: null
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EWallet define payload or parameter for CreateRequest ChannelCode EWallet
type EWallet struct {
	// ReferenceID merchant's identifier for specific transaction id
	// this field required
	//
	// Default: ""
	ReferenceID string `json:"reference_id"`

	// Currency used for the transaction
	// this field required
	//
	// Default: IDR
	Currency Currency `json:"currency"`

	// Amount transaction amount to be paid
	// minimal amount: 100 IDR or 1 PHP
	// maximal based on e-wallet holding limit
	// this field is required
	//
	// Default: 0
	Amount float64 `json:"amount"`

	// CheckoutMethod the payment flow used to be process the transaction
	// this field is required
	//
	// Default: OneTimePayment
	CheckoutMethod CheckoutMethod `json:"checkout_method"`

	// ChannelCode which e-wallet want to be used to process the transaction
	// this field is required
	//
	// Default: ""
	ChannelCode ChannelCode `json:"channel_code"`

	// ChannelProperties any information required for the transaction
	// this filed is required based on CheckoutMethod and ChannelCode pairing
	//
	// Default: null
	ChannelProperties *EWalletChannelProperties `json:"channel_properties"`

	// PaymentMethodID id of PaymentMethod
	// used for tokenized payment
	// this field is required when CheckoutMethod is TokenizedPayment
	//
	// Default: ""
	PaymentMethodID string `json:"payment_method_id,omitempty"`

	// CustomerID id of customer which the payment method will be linked to
	// this field is required when CheckoutMethod is TokenizedPayment
	//
	// Default: ""
	CustomerID string `json:"customer_id,omitempty"`

	// Basket details of items purchased by customer
	//
	// Default: null
	Basket []*EWalletItem `json:"basket,omitempty"`

	// Metadata additional information
	// define the json properties and values based on your information
	//
	// Default: null
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ChargeAction redirection actions to be taken when is_redirect_required return in response is `true`.
type ChargeAction struct {
	// DesktopWebCheckoutURL e-wallet issuer generated URL for web checkout on devices
	// with a stand-alone screen
	DesktopWebCheckoutURL string `json:"desktop_web_checkout_url,omitempty"`

	// MobileWebCheckoutURL e-wallet issuer generated URL for web checkout on mobile devices
	MobileWebCheckoutURL string `json:"mobile_web_checkout_url,omitempty"`

	// MobileDeeplinkCheckoutURL e-wallet issuer generated URL for deeplink checkout on mobile devices
	// jumps directly into e-wallet app for payment confirmation
	MobileDeeplinkCheckoutURL string `json:"mobile_deeplink_checkout_url,omitempty"`

	// QRCheckoutString e-wallet issuer generated qr string for checkout
	// usually on devices with a stand-alone screen
	QRCheckoutString string `json:"qr_checkout_string,omitempty"`
}

// ChargeResponse charge response and notification
type ChargeResponse struct {
	// ID unique identifier for charge request transaction
	ID string `json:"id"`

	// BusinessID of the merchant
	BusinessID string `json:"business_id"`

	// ReferenceID provided by merchant
	ReferenceID string `json:"reference_id"`

	// Status of charge Request
	// possible values:
	// - SUCCEEDED = payment transaction is successfully
	// - PENDING = payment transaction is awaiting payment attempt by end user
	// - FAILED = payment transaction has failed or expired
	// - VOIDED = payment transaction has been voided
	// - REFUNDED = payment transaction has been either partially or fully refunded to end user
	Status ChargeStatus `json:"status"`

	// Currency used for the transaction
	Currency Currency `json:"currency"`

	// ChargeAmount requested charge amount from merchant
	ChargeAmount float64 `json:"charge_amount"`

	// CaptureAmount requested capture amount from merchant
	CaptureAmount float64 `json:"capture_amount,omitempty"`

	// RefundedAmount total amount refunded by merchant to end user
	RefundedAmount float64 `json:"refunded_amount,omitempty"`

	// CheckoutMethod determines the payment flow used to process the transaction
	CheckoutMethod CheckoutMethod `json:"checkout_method"`

	// ChannelCode specifies which e-wallet bill be used to process the transaction
	ChannelCode ChannelCode `json:"channel_code"`

	// ChannelProperties specific information required for the transaction to be initiated
	ChannelProperties *EWalletChannelProperties `json:"channel_properties,omitempty"`

	// Actions redirection into e-wallet apps for complete payment
	Actions *ChargeAction `json:"actions,omitempty"`

	// IsRedirectRequired flag which indicates whether redirection is required
	// for end user to complete payment
	IsRedirectRequired bool `json:"is_redirect_required"`

	// CallbackURL which payment notifications will be sent
	CallbackURL string `json:"callback_url"`

	// Created ISO8601 timestamp for charge object creation
	Created time.Time `json:"created"`

	// Updated ISO8601 timestamp for charge object update
	Updated time.Time `json:"updated"`

	// VoidStatus of the void request
	VoidStatus ChargeStatus `json:"void_status"`

	// VoidedAt timestamp when transaction was voided
	VoidedAt time.Time `json:"voided_at"`

	// CustomerID to will be linked to the transaction
	CustomerID string `json:"customer_id"`

	// PaymentMethodID for end user payment tokens bind with merchant
	// only support with ChannelCode tokenized payments
	PaymentMethodID string `json:"payment_method_id"`

	// FailureCode by end user or e-wallet issuer
	// the failure_code is notified to the merchant
	// in the payment callback
	FailureCode string `json:"failure_code"`

	// ErrorCode by end user or e-wallet issuer
	// the failure_code is notified to the merchant
	// in the payment callback
	ErrorCode string `json:"error_code"`

	// Message error description
	Message string `json:"message"`

	// Basket details of items purchased by customer
	Basket []*EWalletItem `json:"basket"`

	// Metadata additional information
	// define the json properties and values based on your information
	Metadata map[string]interface{} `json:"metadata"`
}

// VirtualAccount response for created payment channel Virtual Account
type VirtualAccount struct {
	// ID unique ID from Xendit for the virtual account
	// use this ID for support escalation and reconciliation
	ID string `json:"id"`

	// ExternalID an ID of your choice which you provide
	ExternalID string `json:"external_id"`

	// OwnerID your xendit BusinessID
	OwnerID string `json:"owner_id"`

	// BankCode bank code for the relevant bank
	BankCode BankCode `json:"bank_code"`

	// MerchantCode prefix for the Virtual Account
	// this is xendit merchant code for aggregator
	MerchantCode string `json:"merchant_code"`

	// AccountNumber complete virtual account number
	AccountNumber string `json:"account_number"`

	// Name of virtual account
	Name string `json:"name"`

	// Currency in which the virtual account operates
	// all bank have IDR currency expect DBS has USD currency
	Currency Currency `json:"currency"`

	// IsSingleUse there are 2 types of va
	// - true: va will be inactive automatically when payment is success
	// - false: remain active when payment is success and can continue to
	//         receive payment using the same VA
	IsSingleUse bool `json:"is_single_use"`

	// IsClosed there are 2 types of va
	// - true: means your customer can only pay amount specified by you.
	//        payment will reject if attempted payment amount deviates from the amount you specified
	// - false: means your customer can pay any amount to the Virtual Account
	IsClosed bool `json:"is_closed"`

	// ExpectedAmount required amount to be paid by your customer for Closed VA (is_closed: true)
	ExpectedAmount float64 `json:"expected_amount"`

	// SuggestedAmount for the VA. only supported for Mandiri and BRI VA
	SuggestedAmount float64 `json:"suggested_amount"`

	// ExpirationDate timestamp of VA expiration time
	// Timezone UTC+0
	ExpirationDate time.Time `json:"expiration_date"`

	// Description of the VA that will be displayed in payment interface
	Description string `json:"description"`

	// Status of VA
	Status StatusVA `json:"status"`

	// ErrorCode by end user or e-wallet issuer
	// the failure_code is notified to the merchant
	// in the payment callback
	ErrorCode string `json:"error_code"`

	// Message error description
	Message string `json:"message"`
}

// VirtualAccountPayment params body from xendit to our callback
// it will use when Customer already successful payment their VA
type VirtualAccountPayment struct {
	// ID unique ID from Xendit for the virtual account
	// use this ID for support escalation and reconciliation
	ID string `json:"id"`

	// PaymentID xendit internal system payment id
	PaymentID string `json:"payment_id"`

	// CallbackVirtualAccountID id of the VA that was paid
	CallbackVirtualAccountID string `json:"callback_virtual_account_id"`

	// ExternalID an ID of your choice which you provide
	ExternalID string `json:"external_id"`

	// BankCode bank code for the relevant bank
	BankCode BankCode `json:"bank_code"`

	// MerchantCode prefix for the Virtual Account
	// this is xendit merchant code for aggregator
	MerchantCode string `json:"merchant_code"`

	// AccountNumber complete virtual account number
	AccountNumber string `json:"account_number"`

	// Amount that was paid to the VA
	Amount float64 `json:"amount"`

	// SenderName of the end user name that paid into the VA
	SenderName string `json:"sender_name"`

	// TransactionTimestamp timestamp of VA expiration time
	// Timezone UTC+0
	TransactionTimestamp time.Time `json:"transaction_timestamp"`

	// PaymentDetail additional information from the Bank
	PaymentDetail struct {
		// Remark that is inputted by the payer when they are about to make a payment
		Remark string `json:"remark"`
	} `json:"payment_detail"`
}

// CreateVirtualAccountParam define payload or parameter for CreateRequest Payment Channel Virtual Account Bank
type CreateVirtualAccountParam struct {
	// ExternalID an ID of your choice which you provide
	//
	// Default: ""
	ExternalID string `json:"external_id"`

	// BankCode bank code for the relevant bank
	//
	// Default: ""
	BankCode BankCode `json:"bank_code"`

	// Name of virtual account
	//
	// Default: ""
	Name string `json:"name"`

	// VirtualAccountNumber complete virtual account number
	//
	// Default: ""
	VirtualAccountNumber string `json:"virtual_account_number,omitempty"`

	// IsSingleUse there are 2 types of va
	// - true: va will be inactive automatically when payment is success
	// - false: remain active when payment is success and can continue to
	//         receive payment using the same VA
	//
	// Default: false
	IsSingleUse bool `json:"is_single_use"`

	// IsClosed there are 2 types of va
	// - true: means your customer can only pay amount specified by you.
	//        payment will reject if attempted payment amount deviates from the amount you specified
	// - false: means your customer can pay any amount to the Virtual Account
	//
	// Default: false
	IsClosed bool `json:"is_closed"`

	// ExpectedAmount required amount to be paid by your customer for Closed VA (is_closed: true)
	// there has minimum payment on each bank
	//
	// For BankMandiri, BankBNI, BankBJB, BankBRI, BankBSI, BankSahabatSampoerna:
	//    minimum payment: Rp 1
	//    maximum payment: Rp 50.000.000.000
	// For BankPermata:
	//    minimum payment: Rp 1
	//    maximum payment: Rp 9.999.999.999
	// For BankBCA:
	//    minimum payment: Rp 10.000
	//    maximum payment: Rp 999.999.999.999
	// For BankDBS:
	//    minimum payment: USD 1
	//    maximum payment: USD 50,000,000,000
	//
	// Default: 0
	ExpectedAmount float64 `json:"expected_amount,omitempty"`

	// SuggestedAmount for the VA. only supported for Mandiri and BRI VA
	//
	// Default: 0
	SuggestedAmount float64 `json:"suggested_amount,omitempty"`

	// ExpirationDate timestamp of VA expiration time
	// Timezone UTC+0
	//
	// Default: +31 years from creation date (request create virtual account)
	ExpirationDate time.Time `json:"expiration_date,omitempty"`

	// Description of the VA that will be displayed in payment interface
	//
	// Default: ""
	Description string `json:"description,omitempty"`
}

// UpdateVirtualAccountParam define payload or parameter for UpdateRequest Virtual Account created Before
type UpdateVirtualAccountParam struct {
	// ID of the virtual account created before
	//
	// Default: ""
	ID string `json:"-"`

	// ExternalID an ID of your choice which you provide
	//
	// Default: ""
	ExternalID string `json:"external_id,omitempty"`

	// IsSingleUse there are 2 types of va
	// - true: va will be inactive automatically when payment is success
	// - false: remain active when payment is success and can continue to
	//         receive payment using the same VA
	//
	// Default: false
	IsSingleUse bool `json:"is_single_use,omitempty"`

	// ExpectedAmount required amount to be paid by your customer for Closed VA (is_closed: true)
	// there has minimum payment on each bank
	//
	// For BankMandiri, BankBNI, BankBJB, BankBRI, BankBSI, BankSahabatSampoerna:
	//    minimum payment: Rp 1
	//    maximum payment: Rp 50.000.000.000
	// For BankPermata:
	//    minimum payment: Rp 1
	//    maximum payment: Rp 9.999.999.999
	// For BankBCA:
	//    minimum payment: Rp 10.000
	//    maximum payment: Rp 999.999.999.999
	// For BankDBS:
	//    minimum payment: USD 1
	//    maximum payment: USD 50,000,000,000
	//
	// Default: 0
	ExpectedAmount float64 `json:"expected_amount,omitempty"`

	// SuggestedAmount for the VA. only supported for Mandiri and BRI VA
	//
	// Default: 0
	SuggestedAmount float64 `json:"suggested_amount,omitempty"`

	// ExpirationDate timestamp of VA expiration time
	// Timezone UTC+0
	//
	// Default: +31 years from creation date (request create virtual account)
	ExpirationDate time.Time `json:"expiration_date"`

	// Description of the VA that will be displayed in payment interface
	//
	// Default: ""
	Description string `json:"description,omitempty"`
}
