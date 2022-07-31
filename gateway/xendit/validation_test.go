package xendit_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
	"github.com/stretchr/testify/assert"
)

func TestValidationParams_ErrorUnimplemented(t *testing.T) {
	var ab []byte

	err := xendit.ValidationParams(ab)

	assert.NotNil(t, err, "error should to be nil")
}
