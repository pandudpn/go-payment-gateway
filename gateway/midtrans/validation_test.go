package midtrans_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
	"github.com/stretchr/testify/assert"
)

func TestValidationParams_ErrorUnimplemented(t *testing.T) {
	var ab []byte

	err := midtrans.ValidationParams(ab)
	assert.NotNil(t, err, "error should not to be nil")
}
