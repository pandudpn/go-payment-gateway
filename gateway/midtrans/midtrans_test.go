package midtrans_test

import (
	"encoding/json"
	"net/http"
	"net/url"

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

func getMockUrlCardTokenSandBox() string {
	return "https://api.sandbox.midtrans.com/v2/token"
}

func getMockUrlRegisterCardSandBox() string {
	return "https://api.sandbox.midtrans.com/v2/card/register"
}

func getMockUrlProduction() string {
	return "https://api.midtrans.com/v2/charge"
}

func getMockUrlCreatePayAccountSandBox() string {
	return "https://api.sandbox.midtrans.com/v2/pay/account"
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

func getMockParamsLinkAccountPay() *midtrans.LinkAccountPay {
	return &midtrans.LinkAccountPay{
		PaymentType: midtrans.PaymentTypeGopay,
		GopayPartner: &midtrans.GopayPartner{
			PhoneNumber: "81212345678",
			CountryCode: "62",
			RedirectURL: "https://www.your-app.com",
		},
	}
}

func getMockParamsLinkAccountPayBytes() []byte {
	b, _ := json.Marshal(getMockParamsLinkAccountPay())

	return b
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
			Bank:     midtrans.BankBCA,
			VANumber: "2122274139",
		},
	}
}

func getMockParamsBankTransferBytes() []byte {
	b, _ := json.Marshal(getMockParamsBankTransfer())
	return b
}

func getMockParamsCardToken() *midtrans.CardToken {
	return &midtrans.CardToken{
		CardNumber:   "5211111111111117",
		CardCvv:      "123",
		CardExpYear:  "2022",
		CardExpMonth: "12",
	}
}

func getMockParamsCardTokenWithSavedTokenId() *midtrans.CardToken {
	return &midtrans.CardToken{
		TokenID: "token-id-testing",
		CardCvv: "123",
	}
}

func getMockParamsRegisterCard() *midtrans.CardRegister {
	return &midtrans.CardRegister{
		CardNumber:   "5211111111111117",
		CardCvv:      "123",
		CardExpYear:  "2022",
		CardExpMonth: "12",
	}
}

func getMockQueryParamsCardToken() string {
	u := url.Values{}
	u.Set("card_cvv", "123")
	u.Set("card_exp_month", "12")
	u.Set("card_exp_year", "2022")
	u.Set("card_number", "5211111111111117")
	u.Set("client_key", clientIdSandBox)

	return u.Encode()
}

func getMockQueryParamsCardTokenBySavedTokenId() string {
	u := url.Values{}
	u.Set("card_cvv", "123")
	u.Set("client_key", clientIdSandBox)
	u.Set("token_id", "token-id-testing")

	return u.Encode()
}

func getMockParamsCardCharge() *midtrans.CardPayment {
	return &midtrans.CardPayment{
		PaymentType: midtrans.PaymentTypeCard,
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
		CreditCard: &midtrans.CreditCard{
			TokenID:        "token-id-testing",
			Authentication: true,
		},
	}
}

func getMockParamsCardChargeBytes() []byte {
	b, _ := json.Marshal(getMockParamsCardCharge())
	return b
}

func getMockAccountId() string {
	return "1ea74cfd-56a0-4e2a-adc4-956e10a04279"
}
