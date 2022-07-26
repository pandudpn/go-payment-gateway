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

	// BankTransferBCA is payment type for Virtual Account with Bank BCA from Midtrans Core API
	// bca from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	BankTransferBCA PaymentType = "bank_transfer"

	// BankTransferBRI is payment type for Virtual Account with Bank BRI from Midtrans Core API
	// bri from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	BankTransferBRI PaymentType = "bank_transfer"

	// BankTransferBNI is payment type for Virtual Account with Bank BNI from Midtrans Core API
	// bni from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	BankTransferBNI PaymentType = "bank_transfer"

	// BankTransferMandiri is payment type for Virtual Account with Bank Mandiri from Midtrans Core API
	// mandiri from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	BankTransferMandiri PaymentType = "echannel"

	// BankTransferPermata is payment type for Virtual Account with Bank Permata from Midtrans Core API
	// permata from Midtrans only support transaction OpenAmount (need to fill the amount after input VA Number)
	BankTransferPermata PaymentType = "bank_transfer"
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
)
