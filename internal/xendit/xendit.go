package xendit

import (
	"context"
)

// PaymentInterface define method to be implemented by each PaymentType
type PaymentInterface interface {
	// SetUsername will set username in Basic Auth
	SetUsername(username string)
	
	// SetURI set url target
	SetURI(uri string)
	
	// Do create a charge payments
	Do(ctx context.Context) (*ChargeResponse, error)
}
