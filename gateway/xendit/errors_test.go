package xendit_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
	"github.com/stretchr/testify/assert"
)

func TestGetErrorCode_Found(t *testing.T) {
	errorCode := "CHANNEL_UNAVAILABLE"

	expected := xendit.ErrChannelUnavailable

	f := xendit.GetErrorCode(errorCode)
	assert.Equal(t, expected, f)
}

func TestGetErrorCode_NotFound(t *testing.T) {
	errorCode := "abc"

	expected := xendit.ErrUnknown

	f := xendit.GetErrorCode(errorCode)
	assert.Equal(t, expected, f)
}
