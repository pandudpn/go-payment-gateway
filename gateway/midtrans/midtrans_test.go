package midtrans_test

import (
	"encoding/json"
	"net/http"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
)

const clientIdSandBox string = "sb-client_id"
const clientId string = "client_id"
const clientSecretSandBox string = "sb-client_secret"
const clientSecret string = "client_secret"

func init() {
	pg.NewLogger()
}

func getMockUrlSandBox() string {
	return "https://api.sandbox.midtrans.com/v2/charge"
}

func getMockUrlProduction() string {
	return "https://api.midtrans.com/v2/charge"
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
	header.Set("Authorization", "Basic c2ItY2xpZW50X3NlY3JldDo=")

	return header
}

func getMockHeaderProduction() http.Header {
	var header = make(http.Header)
	header.Set("Authorization", "Basic Y2xpZW50X3NlY3JldDo=")

	return header
}

func getMockParamsEWallet() *midtrans.EWallet {
	return &midtrans.EWallet{
		PaymentType: midtrans.PaymentTypeGopay,
		TransactionDetails: &midtrans.TransactionDetail{
			OrderID:     "order-id",
			GrossAmount: 10000,
		},
		ItemDetails: []*midtrans.ItemDetail{
			{
				ID:       "item-id",
				Price:    10000,
				Name:     "abc",
				Quantity: 1,
			},
		},
	}
}

func getMockParamsEWalletBytes() []byte {
	b, _ := json.Marshal(getMockParamsEWallet())

	return b
}

func getMockChargeResponseEWallet() *midtrans.ChargeResponse {
	return &midtrans.ChargeResponse{
		ID:                "id-test",
		StatusCode:        "201",
		StatusMessage:     "GoPay transaction is created",
		TransactionID:     "transaction-id-test",
		OrderID:           "order-id-test",
		GrossAmount:       "10000",
		PaymentType:       midtrans.PaymentTypeGopay,
		TransactionStatus: "settlement",
		FraudStatus:       "accept",
		TransactionTime:   "2022-07-30 10:00:00",
		Actions: []*midtrans.Action{
			{
				Name:   "deeplink-redirect",
				Method: "GET",
				URL:    "https://simulator.sandbox.midtrans.com/gopay/partner/app/payment-pin?id=fe51909d-cf14-42ff-af57-788fe97a74e3",
			},
		},
	}
}

func getMockParamsBankTransfer() *midtrans.BankTransferCreateParams {
	return &midtrans.BankTransferCreateParams{
		PaymentType: midtrans.PaymentTypeBCA,
		TransactionDetails: &midtrans.TransactionDetail{
			OrderID:     "order-id-test",
			GrossAmount: 10000,
		},
		ItemDetails: []*midtrans.ItemDetail{
			{
				ID:       "item-id",
				Price:    10000,
				Name:     "abc",
				Quantity: 1,
			},
		},
		BankTransfer: &midtrans.BankTransfer{
			Bank:     midtrans.BankTransferBCA,
			VANumber: "2122274139",
		},
	}
}

func getMockParamsBankTransferBytes() []byte {
	b, _ := json.Marshal(getMockParamsBankTransfer())
	return b
}

func getMockChargeResponseBankTransfer() *midtrans.ChargeResponse {
	return &midtrans.ChargeResponse{
		ID:                "id-test",
		StatusCode:        "201",
		StatusMessage:     "Success, Bank Transfer transaction is created",
		TransactionID:     "transaction-id-test",
		OrderID:           "order-id-test",
		GrossAmount:       "10000",
		PaymentType:       midtrans.PaymentTypeBCA,
		TransactionStatus: "pending",
		FraudStatus:       "accept",
		TransactionTime:   "2022-07-30 10:00:00",
		VANumbers: []*midtrans.BankTransfer{
			{
				Bank:     midtrans.BankTransferBCA,
				VANumber: "268192122274139",
			},
		},
	}
}
