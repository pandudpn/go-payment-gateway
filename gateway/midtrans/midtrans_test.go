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

func getMockUrlCardTokenProduction() string {
	return "https://api.midtrans.com/v2/token"
}

func getMockUrlRegisterCardProduction() string {
	return "https://api.midtrans.com/v2/card/register"
}

func getMockUrlCreatePayAccountSandBox() string {
	return "https://api.sandbox.midtrans.com/v2/pay/account"
}

func getMockUrlCreatePayAccountProduction() string {
	return "https://api.midtrans.com/v2/pay/account"
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

func getMockLinkAccountPayResponse() *midtrans.LinkAccountPayResponse {
	return &midtrans.LinkAccountPayResponse{
		PaymentType: midtrans.PaymentTypeGopay,
		StatusCode:  "201",
		AccountID:   "account-id-testing",
		Actions: []*midtrans.Action{
			{
				Name:   "activation-deeplink",
				Method: "GET",
				URL:    "https://api.sandbox.midtrans.com/v2/pay/account/gpar_account-id-testing/link",
			},
		},
		AccountStatus: "PENDING",
	}
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
			Bank:     midtrans.BankBCA,
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
				Bank:     midtrans.BankBCA,
				VANumber: "268192122274139",
			},
		},
	}
}

func getMockParamsCardToken() *midtrans.CardToken {
	return &midtrans.CardToken{
		CardNumber:   "5211111111111117",
		CardCvv:      "123",
		CardExpYear:  "2022",
		CardExpMonth: "12",
	}
}

func getMockParamsCardRegisterBytes() []byte {
	b, _ := json.Marshal(getMockParamsCardToken())
	return b
}

func getMockParamsCardTokenWithSavedTokenId() *midtrans.CardToken {
	return &midtrans.CardToken{
		TokenID: "token-id-testing",
		CardCvv: "123",
	}
}

func getMockParamsCardTokenWithSavedTokenIdBytes() []byte {
	b, _ := json.Marshal(getMockParamsCardTokenWithSavedTokenId())
	return b
}

func getMockParamsRegisterCard() *midtrans.CardRegister {
	return &midtrans.CardRegister{
		CardNumber:   "5211111111111117",
		CardCvv:      "123",
		CardExpYear:  "2022",
		CardExpMonth: "12",
	}
}

func getMockParamsRegisterCardBytes() []byte {
	b, _ := json.Marshal(getMockParamsRegisterCard())
	return b
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

func getMockQueryParamsCardTokenBytes() []byte {
	b, _ := json.Marshal(getMockQueryParamsCardToken())
	return b
}

func getMockQueryParamsCardTokenBySavedTokenId() string {
	u := url.Values{}
	u.Set("card_cvv", "123")
	u.Set("client_key", clientIdSandBox)
	u.Set("token_id", "token-id-testing")

	return u.Encode()
}

func getMockQueryParamsCardTokenBySavedTokenIdBytes() []byte {
	b, _ := json.Marshal(getMockQueryParamsCardTokenBySavedTokenId())
	return b
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
