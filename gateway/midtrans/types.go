package midtrans

// TransactionDetail for customer
type TransactionDetail struct {
	// OrderID is Reference ID for midtrans
	//
	// Default: uuid v4
	OrderID string `json:"order_id"`

	// GrossAmount is total transaction to be paid by Customer
	// if value decimal, it will round up
	GrossAmount int64 `json:"gross_amount"`
}

// CustomerDetail details of customer
type CustomerDetail struct {
	// FirstName customer first name
	FirstName string `json:"first_name,omitempty"`

	// LastName customer last name
	LastName string `json:"last_name,omitempty"`

	// Email customer email address
	Email string `json:"email,omitempty"`

	// Phone customer phone number
	Phone string `json:"phone,omitempty"`

	// BillingAddress customer billing address
	BillingAddress string `json:"billing_address,omitempty"`

	// ShippingAddress customer shipping address
	ShippingAddress string `json:"shipping_address"`
}

// ItemDetail details items purchased by Customer
type ItemDetail struct {
	// ID product_id from catalogue
	ID string `json:"id,omitempty"`

	// Name product_name from catalogue
	Name string `json:"name"`

	// Price product price (after discount) from catalogue
	Price int64 `json:"price"`

	// Quantity of the item purchased by customer
	Quantity int64 `json:"quantity"`

	// Brand product brand from catalogue
	Brand string `json:"brand,omitempty"`

	// Category product category from catalogue
	Category string `json:"category,omitempty"`

	// Merchant name of merchant selling the item
	Merchant string `json:"merchant,omitempty"`

	// URL of item in the merchant site
	URL string `json:"url,omitempty"`
}

// EWalletDetail details e-wallet payment using GoPay or ShopeePay
type EWalletDetail struct {
	// CallbackURL http or deeplink URL to which customer is redirected
	// from Gojek app or Shopee app after successful payment
	//
	// Default: from dashboard settings
	CallbackURL string `json:"callback_url,omitempty"`

	// EnableCallback required for gopay deeplink
	// to determine appending callback url in deeplink
	//
	// Default: false
	EnableCallback bool `json:"enable_callback,omitempty"`

	// AccountID required for gopay tokenization (linked account)
	AccountID string `json:"account_id,omitempty"`

	// PaymentOptionToken required for gopay tokenization (linked account)
	// token to specify the payment option GOPAY_WALLET or PAY_LATER
	PaymentOptionToken string `json:"payment_option_token,omitempty"`

	// PreAuth set the value to TRUE to reserve the specified amount
	// from the customer balance
	//
	// Default: false
	PreAuth bool `json:"pre_auth,omitempty"`

	// Recurring set the value to TRUE to mark as a recurring transaction
	//
	// Default: false
	Recurring bool `json:"recurring,omitempty"`
}

// EWallet define payload or parameter for request CreatePayment EWallet
type EWallet struct {
	// PaymentType is payment method
	PaymentType PaymentType `json:"payment_type"`

	// TransactionDetails the details of the specific transactions
	TransactionDetails *TransactionDetail `json:"transaction_details"`

	// ItemDetails details of items purchased by customer
	ItemDetails []*ItemDetail `json:"item_details"`

	// CustomerDetails detail of customer
	CustomerDetails *CustomerDetail `json:"customer_details,omitempty"`

	// Gopay charge detail using Gopay (Tokenization or OneTimePayment)
	Gopay *EWalletDetail `json:"gopay,omitempty"`

	// ShopeePay charge detail using ShopeePay (OneTimePayment)
	ShopeePay *EWalletDetail `json:"shopeepay,omitempty"`
}

// Action to make payments redirect
type Action struct {
	// Name action name
	Name string `json:"name"`

	// Method http method used for the action
	Method string `json:"method"`

	// URL http url target for the action
	URL string `json:"url"`

	// Fields parameters which can be sent for the action
	Fields []string `json:"fields"`
}

// EChannel charge details using Mandiri Bill Payment
type EChannel struct {
	// BillInfo1 label 1 allows only 10 characters
	// this is required
	//
	// Default: ""
	BillInfo1 string `json:"bill_info1"`

	// BillInfo2 value for label 1. allows only 30 characters
	// this is required
	//
	// Default: ""
	BillInfo2 string `json:"bill_info2"`

	// BillInfo3 label 2 allows only 10 characters
	//
	// Default: ""
	BillInfo3 string `json:"bill_info3,omitempty"`

	// BillInfo4 value for label 2. allows only 30 characters
	//
	// Default: ""
	BillInfo4 string `json:"bill_info4,omitempty"`

	// BillInfo5 label 3 allows only 10 characters
	//
	// Default: ""
	BillInfo5 string `json:"bill_info5,omitempty"`

	// BillInfo6 value for label 3. allows only 30 characters
	//
	// Default: ""
	BillInfo6 string `json:"bill_info6,omitempty"`

	// BillInfo7 label 4 allows only 10 characters
	//
	// Default: ""
	BillInfo7 string `json:"bill_info7,omitempty"`

	// BillInfo8 value for label 4. allows only 30 characters
	//
	// Default: ""
	BillInfo8 string `json:"bill_info8,omitempty"`

	// BillKey custom VA Number
	// notes:: if the previous transaction with same bill key
	// is still active, then Midtrans will return new random bill key
	BillKey string `json:"bill_key,omitempty"`
}

// BCA specific parameters for BCA VA
type BCA struct {
	// SubCompanyCode BCA sub company code directed for this transactions
	//
	// Default: "00000"
	SubCompanyCode string `json:"sub_company_code"`
}

// Permata specific parameters for Permata VA
type Permata struct {
	// RecipientName shown on the payment details
	//
	// Default: Merchant Name. e.g: PanduShop
	RecipientName string `json:"recipient_name"`
}

// FreeTextMessage for message in language
type FreeTextMessage struct {
	// ID free text message in Bahasa Indonesia
	ID string `json:"id"`

	// EN free text message in English
	EN string `json:"en"`
}

// FreeText list of free texts used for BCA VA
type FreeText struct {
	// Inquiry free texts shown during inquiry
	Inquiry []*FreeTextMessage `json:"inquiry"`

	// Payment free texts shown during payment
	Payment []*FreeTextMessage `json:"payment"`
}

// BankTransfer charge details using bank transfer
type BankTransfer struct {
	// Bank name which processes bank transfer transaction
	// this field is required
	//
	// Default: ""
	Bank BankCode `json:"bank"`

	// VANumber for customized virtual_account number
	// notes::
	// - BCA: accepts 6-11 digits
	// - Permata: accepts 10 digits
	// - BNI: accepts 8-12 digits
	// - BRI: accepts 18 digits
	//
	// Default: ""
	VANumber string `json:"va_number,omitempty"`

	// FreeText list of free text
	//
	// Default: null
	FreeText *FreeText `json:"free_text,omitempty"`

	// BCA specific params for BCA VA
	//
	// Default: null
	BCA *BCA `json:"bca,omitempty"`

	// Permata specific params for Permata VA
	//
	// Default: null
	Permata *Permata `json:"permata,omitempty"`
}

// BankTransferCreateParams parameters to create Bank Transfer transaction
type BankTransferCreateParams struct {
	// PaymentType set Bank Transfer payment method
	//
	// Default: ""
	PaymentType PaymentType `json:"payment_type"`

	// TransactionDetails the details of the specific transactions
	TransactionDetails *TransactionDetail `json:"transaction_details"`

	// ItemDetails details of items purchased by customer
	ItemDetails []*ItemDetail `json:"item_details"`

	// CustomerDetails detail of customer
	CustomerDetails *CustomerDetail `json:"customer_details,omitempty"`

	// BankTransfer charge details using bank transfer
	BankTransfer *BankTransfer `json:"bank_transfer,omitempty"`

	// EChannel charge details using Mandiri Bill Payment
	EChannel *EChannel `json:"echannel,omitempty"`
}

// ChargeResponse charge response and notifications
type ChargeResponse struct {
	// ID of request
	ID string `json:"id"`

	// StatusCode status code of transaction result
	StatusCode string `json:"status_code"`

	// StatusMessage description of transaction result
	StatusMessage string `json:"status_message"`

	// ValidationMessages parameters validation
	ValidationMessages []string `json:"validation_messages"`

	// TransactionID transaction id given by Midtrans
	TransactionID string `json:"transaction_id"`

	// OrderID order id specified by merchant
	OrderID string `json:"order_id"`

	// GrossAmount total amount of transaction in IDR
	GrossAmount string `json:"gross_amount"`

	// PaymentType transaction payment method
	PaymentType PaymentType `json:"payment_type"`

	// TransactionTime timestamp of transaction in ISO 8601 format with timezone GMT+7
	TransactionTime string `json:"transaction_time"`

	// TransactionStatus status of transaction.
	// values:
	// - pending
	// - settlement
	// - expire
	// - deny
	TransactionStatus string `json:"transaction_status"`

	// FraudStatus detection result by Fraud Detection System (FDS)
	// values:
	// - accept : approved by FDS
	// - challenge : questioned by FDS
	// - deny : denied by FDS. transaction automatically failed
	FraudStatus string `json:"fraud_status"`

	// ChannelResponseCode response code from payment channel provider
	ChannelResponseCode string `json:"channel_response_code"`

	// ChannelResponseMessage response message from payment channel provider
	ChannelResponseMessage string `json:"channel_response_message"`

	// Actions to take action payments e-wallet
	Actions []*Action `json:"actions"`

	// BillKey va_number for payment Bank Mandiri
	BillKey string `json:"bill_key"`

	// BillerCode company code for Bank Mandiri
	BillerCode string `json:"biller_code"`

	// PermataVANumber va_number for payment Permata Bank
	PermataVANumber string `json:"permata_va_number"`

	// VANumbers list va_number for payment Bank BCA, BRI, or BNI
	VANumbers []*BankTransfer `json:"va_numbers"`
}
