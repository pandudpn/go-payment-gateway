package xendit

import (
	"encoding/json"
)

func (e *EWallet) CreateRequest() PaymentInterface {
	payload, _ := json.Marshal(e)

	return &request{params: payload}
}
