package xendit_test

import (
	"encoding/json"
	"net/http"
	"time"

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

func getMockCreateVirtualAccount() *xendit.CreateVirtualAccountParam {
	layoutDateTime := "2006-01-02T15:04:05"
	expired, _ := time.Parse(layoutDateTime, "2022-08-17T23:59:59")

	return &xendit.CreateVirtualAccountParam{
		ExternalID:           "external-id-test",
		ExpectedAmount:       10000,
		BankCode:             xendit.BankBCA,
		ExpirationDate:       expired,
		IsClosed:             true,
		IsSingleUse:          true,
		Name:                 "Pandu dwi Putra",
		VirtualAccountNumber: "9999345678",
	}
}

func getMockCreateVirtualAccountBytes() []byte {
	b, _ := json.Marshal(getMockCreateVirtualAccount())

	return b
}

func getMockVirtualAccount() *xendit.VirtualAccount {
	layoutDateTime := "2006-01-02T15:04:05"
	expired, _ := time.Parse(layoutDateTime, "2022-08-17T23:59:59")

	return &xendit.VirtualAccount{
		ID:             "id-virtual-account-testing",
		Name:           "Pandu dwi Putra",
		IsClosed:       true,
		IsSingleUse:    true,
		BankCode:       xendit.BankBCA,
		ExpirationDate: expired,
		ExpectedAmount: 10000,
		Currency:       xendit.IDR,
		ExternalID:     "external-id-test",
		Status:         xendit.PendingVA,
		AccountNumber:  "107669999345678",
		MerchantCode:   "10766",
	}
}

func getMockUpdateVirtualAccount() *xendit.UpdateVirtualAccountParam {
	layoutDateTime := "2006-01-02T15:04:05"
	expired, _ := time.Parse(layoutDateTime, "2022-08-17T23:59:59")

	return &xendit.UpdateVirtualAccountParam{
		ExternalID:     "external-id-test",
		ExpectedAmount: 10000,
		ExpirationDate: expired,
		IsSingleUse:    true,
		ID:             "id-virtual-account-testing",
	}
}

func getMockUpdateVirtualAccountBytes() []byte {
	b, _ := json.Marshal(getMockUpdateVirtualAccount())

	return b
}
