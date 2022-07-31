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

// CustomerAddress address of customer detail
type CustomerAddress struct {
	// FirstName address of customer first name
	FirstName string `json:"first_name,omitempty"`

	// LastName address of customer last name
	LastName string `json:"last_name,omitempty"`

	// Phone customer phone number
	Phone string `json:"phone,omitempty"`

	// Address of customer
	Address string `json:"address,omitempty"`

	// City of customer
	City string `json:"city,omitempty"`

	// PostalCode zip code of customer address
	PostalCode string `json:"postal_code,omitempty"`

	// CountryCode code of country customer
	CountryCode string `json:"country_code"`
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
	BillingAddress *CustomerAddress `json:"billing_address,omitempty"`

	// ShippingAddress customer shipping address
	ShippingAddress *CustomerAddress `json:"shipping_address,omitempty"`
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

// GopayPartner gopay linking specific parameters
type GopayPartner struct {
	// PhoneNumber linked to the customer's account
	//
	// Default: ""
	PhoneNumber string `json:"phone_number"`

	// CountryCode associated with the phone number
	//
	// Default: ""
	CountryCode string `json:"country_code"`

	// RedirectURL where user is redrected to after finishing
	// the confirmation on Gojek app
	//
	// Default: ""
	RedirectURL string `json:"redirect_url"`
}

// PaymentOption additional data from the specific payment provider
type PaymentOption struct {
	// Name payment option
	Name string `json:"name"`

	// Active status of payment options
	Active bool `json:"active"`

	// Token that need to be use on ChargeEWallet as PaymentOptionToken
	Token string `json:"token"`

	// Balance linked account balance for each payment option
	Balance struct {
		// Value balance
		Value string `json:"value"`

		// Currency balance
		Currency string `json:"currency"`
	} `json:"balance"`
}

// LinkAccountPay is triggered to link the customer's account
// to be used for payments using specific payment channel
type LinkAccountPay struct {
	// PaymentType Channel where the account is register to
	//
	// Default: PaymentTypeGopay
	PaymentType PaymentType `json:"payment_type"`

	// GopayPartner GoPay linking params
	//
	// Default: null
	GopayPartner *GopayPartner `json:"gopay_partner"`
}

// LinkAccountPayResponse response from creating Pay Account
type LinkAccountPayResponse struct {
	// StatusCode of the API Result
	StatusCode string `json:"status_code"`

	// PaymentType channel associated with the account
	PaymentType PaymentType `json:"payment_type"`

	// AccountID customer account id to be used for payment
	AccountID string `json:"account_id"`

	// AccountStatus status of the account
	// possible values are:
	// - PENDING
	// - EXPIRED
	// - ENABLED
	// - DISABLED
	AccountStatus string `json:"account_status"`

	// ChannelResponseCode response code from payment channel provider
	ChannelResponseCode string `json:"channel_response_code"`

	// ChannelResponseMessage response message from payment channel provider
	ChannelResponseMessage string `json:"channel_response_message"`

	// Actions to be performed
	Actions []*Action `json:"actions"`

	// Metadata additional data
	Metadata struct {
		// PaymentOptions data from the specific provider payment
		PaymentOptions []*PaymentOption `json:"payment_options"`

		// ReferenceID identifier for specific request
		ReferenceID string `json:"reference_id"`
	} `json:"metadata"`
}

// CardToken create a tokenization for credit_card or debit_card
type CardToken struct {
	// TokenID The token ID of credit card saved previously
	TokenID string `json:"token_id,omitempty"`

	// CardNumber which will be converted into a secured token
	// this field is required
	//
	// Default: ""
	CardNumber string `json:"card_number,omitempty"`

	// CardExpMonth expired month of card
	// this field is required
	//
	// Default: ""
	CardExpMonth string `json:"card_exp_month,omitempty"`

	// CardExpYear expired year of card
	// this field is required
	//
	// Default: ""
	CardExpYear string `json:"card_exp_year,omitempty"`

	// CardCvv three digit unique written on the back card
	// this field is required
	//
	// Default: ""
	CardCvv string `json:"card_cvv"`
}

// CardRegister create a tokenization for credit_card or debit_card
// and will save the token for future transactions
// **notes : no need more input card_number, card_exp_month, or card_exp_year
//			 after Register the Card
type CardRegister struct {
	// CardNumber which will be converted into a secured token
	// this field is required
	//
	// Default: ""
	CardNumber string `json:"card_number,omitempty"`

	// CardExpMonth expired month of card
	// this field is required
	//
	// Default: ""
	CardExpMonth string `json:"card_exp_month,omitempty"`

	// CardExpYear expired year of card
	// this field is required
	//
	// Default: ""
	CardExpYear string `json:"card_exp_year,omitempty"`

	// CardCvv three digit unique written on the back card
	// this field is required
	//
	// Default: ""
	CardCvv string `json:"card_cvv"`
}

// CreditCard the details of payment used for the transaction
type CreditCard struct {
	// TokenID represents customer credit card information
	//
	// Default: ""
	TokenID string `json:"token_id"`

	// Bank Name of the acquiring bank
	//
	// Default: ""
	Bank BankCode `json:"bank,omitempty"`

	// InstallmentTerm tenure in terms of months
	//
	// Default: 0
	InstallmentTerm int `json:"installment_term,omitempty"`

	// Bins List of credit card's BIN (Bank Identification Number)
	// that is allowed for transaction
	//
	// Default: null
	Bins []string `json:"bins,omitempty"`

	// Type Used as preAuthorization feature
	// valid value
	// - authorize
	//
	// Default: ""
	Type string `json:"type,omitempty"`

	// Authentication Flag to enable the 3D secure authentication
	//
	// Default: false
	Authentication bool `json:"authentication,omitempty"`

	// SaveTokenID Used on 'One Click' or 'Two Clicks' feature
	// Enabling it will return a saved_token_id
	// that can be used for the next transaction
	//
	// Default: false
	SaveTokenID bool `json:"save_token_id,omitempty"`
}

// CardPayment params charge transaction using credit_card or debit_card
type CardPayment struct {
	// PaymentType set Bank Transfer payment method
	//
	// Default: PaymentTypeCard
	PaymentType PaymentType `json:"payment_type"`

	// TransactionDetails the details of the specific transactions
	//
	// Default: null
	TransactionDetails *TransactionDetail `json:"transaction_details"`

	// ItemDetails details of items purchased by customer
	//
	// Default: null
	ItemDetails []*ItemDetail `json:"item_details"`

	// CustomerDetails detail of customer
	//
	// Default: null
	CustomerDetails *CustomerDetail `json:"customer_details,omitempty"`

	// CreditCard the details of payment used for the transaction
	//
	// Default: null
	CreditCard *CreditCard `json:"credit_card"`
}

// CardResponse after create or register CreditCard
type CardResponse struct {
	// StatusCode status code of transaction result
	StatusCode string `json:"status_code"`

	// StatusMessage description of transaction result
	StatusMessage string `json:"status_message"`

	// TokenID token transaction for payment CardToken
	TokenID string `json:"token_id"`

	// SavedTokenID A flag to indicate whether the token_id is saved for future transactions
	SavedTokenID string `json:"saved_token_id"`

	// TransactionID transaction id given by Midtrans
	TransactionID string `json:"transaction_id"`

	// MaskedCard first 6-digit and last 4-digit of customer's payment card number
	MaskedCard string `json:"masked_card"`

	// Hash algorithm for hashing CardToken
	Hash string `json:"hash"`
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

	// SavedTokenID A flag to indicate whether the token_id is saved for future transactions
	SavedTokenID string `json:"saved_token_id"`

	// RedirectURL The URL to redirect the user to 3D Secure authentication page
	RedirectURL string `json:"redirect_url"`

	// Bank code
	Bank BankCode `json:"bank"`

	// MaskedCard first 6-digit and last 4-digit of customer's payment card number
	MaskedCard string `json:"masked_card"`

	// CardType type of card used
	// possible values:
	// - credit
	// - debit
	CardType string `json:"card_type"`
}
