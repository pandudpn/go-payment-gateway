package pg_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway"
	"github.com/stretchr/testify/assert"
)

func TestNewOption_Success(t *testing.T) {
	var err error

	expectedResult := getMockOptionsFalse()
	opts := &pg.Options{
		ServerKey: "abc",
		ClientId:  "abc",
		Logging:   &pg.False,
	}

	opts, err = pg.NewOption(opts)
	expectedResult.ApiCall = opts.ApiCall

	assert.Equal(t, expectedResult, opts)
	assert.Nil(t, err, "error should be nil")
}

func TestNewOption_ErrorMissingCredentials(t *testing.T) {
	opts, err := pg.NewOption()

	assert.Nil(t, opts, "options should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}
