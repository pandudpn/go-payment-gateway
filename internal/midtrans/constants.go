package midtrans

type PaymentType string

// list midtrans payment type
const (
	// EWalletGopay is payment type gopay from Midtrans Core API
	// gopay supporting for payment:
	//  - One Time Payment
	//  - Linked Account (Tokenized)
	//  - Recurring
	EWalletGopay PaymentType = "gopay"
	
	// EWalletShopeePay is payment type ShopeePay from Midtrans Core API
	// shopeePay from Midtrans only support One Time Payment
	EWalletShopeePay PaymentType = "shopeepay"
)

// convert to String value or type
func (pt PaymentType) String() string {
	return string(pt)
}
