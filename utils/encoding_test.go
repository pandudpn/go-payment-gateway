package utils_test

import (
	"testing"

	"github.com/pandudpn/go-payment-gateway/utils"
)

func TestSetBasicAuthorization(t *testing.T) {
	utils.SetBasicAuthorization("", "")
}
