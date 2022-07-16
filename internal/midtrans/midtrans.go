package midtrans

// PaymentInterface define method to be implemented by each PaymentType
type PaymentInterface interface {
	// CreateRequest will create a new instance of Charge Payment
	CreateRequest() *request
}
