package xendit_test

import (
	"encoding/json"
	"net/http"

	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
)

const clientIdSandBox string = "xnd_public_development_testing"
const clientId string = "xnd_public_production_testing"
const clientSecretSandBox string = "xnd_development_testing"
const clientSecret string = "xnd_production_testing"

func init() {
	pg.NewLogger()
}

func getMockUrlEWallet() string {
	return "https://api.xendit.co/ewallets/charges"
}

func getMockOptionsSandBox() *pg.Options {
	return &pg.Options{
		Environment: pg.SandBox,
		ClientId:    clientIdSandBox,
		ServerKey:   clientSecretSandBox,
		Logging:     &pg.True,
	}
}

func getMockOptionsProduction() *pg.Options {
	return &pg.Options{
		Environment: pg.Production,
		ClientId:    clientId,
		ServerKey:   clientSecret,
		Logging:     &pg.True,
	}
}

func getMockHeaderSandBox() http.Header {
	var header = make(http.Header)
	header.Set("Authorization", "Basic eG5kX2RldmVsb3BtZW50X3Rlc3Rpbmc6")

	return header
}

func getMockHeaderProduction() http.Header {
	var header = make(http.Header)
	header.Set("Authorization", "Basic eG5kX3Byb2R1Y3Rpb25fdGVzdGluZzo=")

	return header
}

func getMockParamsEWallet() *xendit.EWallet {
	return &xendit.EWallet{
		ChannelCode:    xendit.ChannelCodeShopeePay,
		Amount:         10000,
		Currency:       xendit.IDR,
		CheckoutMethod: xendit.OneTimePayment,
		ReferenceID:    "ref-id-testing",
		ChannelProperties: &xendit.EWalletChannelProperties{
			SuccessRedirectURL: "https://www.google.com",
		},
	}
}

func getMockParamsEWalletBytes() []byte {
	b, _ := json.Marshal(getMockParamsEWallet())

	return b
}

func getMockChargeResponseEWallet() *xendit.ChargeResponse {
	return &xendit.ChargeResponse{
		ID: "id-test",
	}
}
