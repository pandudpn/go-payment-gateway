package xendit

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pandudpn/go-payment-gateway/utils"
)

// SetUsername for Authorization Basic
func (r *request) SetUsername(username string) {
	r.username = username
}

// SetURI target http api
func (r *request) SetURI(uri string) {
	r.uri = uri
}

// Do create a charge payments
func (r *request) Do(ctx context.Context) (*ChargeResponse, error) {
	// create instance utils.Request
	req, err := utils.NewRequest(http.MethodPost, r.uri, r.params)
	if err != nil {
		return nil, err
	}
	// create a header basic authorization
	req.SetBasicAuth(r.username, "")

	// request to target
	// return response body with array bytes
	// and http.StatusCode
	res, statusCode, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}

	var charge *ChargeResponse
	err = json.Unmarshal(res, &charge)
	if err != nil {
		utils.Log.Error("[unmarshal] %s", err)
		return nil, err
	}

	// given error when status_code from header or status_code from body response more than 400
	if statusCode >= http.StatusBadRequest {
		utils.Log.Error(statusCode, charge.Message, charge.ErrorCode)
		return nil, GetErrorCode(charge.ErrorCode)
	}

	return charge, nil
}
