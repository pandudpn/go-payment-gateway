package midtrans

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
)

// convert to String value or type
func (pt PaymentType) String() string {
	return string(pt)
}

type BankCode string

const (
	// BankTransferBCA code for Bank BCA
	BankTransferBCA BankCode = "bca"

	// BankTransferBRI code for Bank BRI
	BankTransferBRI BankCode = "bri"

	// BankTransferBNI code for Bank BNI
	BankTransferBNI BankCode = "bni"

	// BankTransferPermata code for Bank Permata
	BankTransferPermata BankCode = "permata"
)
