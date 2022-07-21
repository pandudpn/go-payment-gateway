package midtrans

import (
	"encoding/json"
)

// CreateRequest will create a new instance of Charge Payment EWallet
func (e *EWallet) CreateRequest() PaymentInterface {
	payload, _ := json.Marshal(e)

	return &request{params: payload}
}
