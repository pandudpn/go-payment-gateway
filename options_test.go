package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOption_Success(t *testing.T) {
	var err error

	expectedResult := getMockOptionsFalse()
	opts := &Options{
		ServerKey: "abc",
		ClientId:  "abc",
		Logging:   &False,
	}

	opts, err = NewOption(opts)
	expectedResult.ApiCall = opts.ApiCall

	assert.Equal(t, expectedResult, opts)
	assert.Nil(t, err, "error should be nil")
}

func TestNewOption_ErrorMissingCredentials(t *testing.T) {
	opts, err := NewOption()

	assert.Nil(t, opts, "options should be nil")
	assert.NotNil(t, err, "error should not be nil")
}
