package xendit_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
	"github.com/stretchr/testify/assert"
)

func TestGetErrorCode_Found(t *testing.T) {
	mockChargeRes := getMockChargeResponseEWallet()
	mockChargeRes.ErrorCode = "API_VALIDATION_ERROR"

	expected := xendit.ErrAPIValidation

	f := xendit.GetErrorCode(*mockChargeRes)
	assert.Equal(t, expected, f)
}

func TestGetErrorCode_NotFound(t *testing.T) {
	mockChargeRes := getMockChargeResponseEWallet()
	mockChargeRes.ErrorCode = "API"

	expected := xendit.ErrUnknown

	f := xendit.GetErrorCode(*mockChargeRes)
	assert.Equal(t, expected, f)
}
