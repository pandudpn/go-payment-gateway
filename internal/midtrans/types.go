package midtrans

type request struct {
	// uri http target to charge payments
	uri string

	// username for authentication basic
	username string

	// params request to charge payment
	params []byte
}

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
	PaymentType string `json:"payment_type"`

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

	// Actions to take action payments
	Actions []*Action `json:"actions"`
}
