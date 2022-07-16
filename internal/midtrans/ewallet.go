package midtrans

import (
	"encoding/json"
)

func (e *EWallet) CreateRequest() *request {
	payload, _ := json.Marshal(e)

	return &request{params: payload}
}
