package midtrans

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strings"

	"github.com/pandudpn/go-payment-gateway"
)

const (
	// uri for production environment
	productionUri = "https://api.midtrans.com"

	// uri for development environment (staging or sandBox)
	sandBoxUri = "https://api.sandbox.midtrans.com"

	// uri for charge transaction
	chargeUri = "/v2/charge"

	// uri for create CardToken
	createCardTokenUri = "/v2/token"

	// uri for register CardToken
	createRegisterCardUri = "/v2/card/register"

	// uri for createPayAccount
	createPayAccountUri = "/v2/pay/account"
)

type midtrans struct {
	// opts configuration or credentials for http api call
	opts *pg.Options

	// uri http target to charge payments
	uri string

	// params request to charge payment
	params []byte
}

func createChargeMidtrans(params interface{}, opts *pg.Options) (*midtrans, error) {
	// return an error when params not exists
	if params == nil || (reflect.ValueOf(params).Kind() == reflect.Ptr && reflect.ValueOf(params).IsNil()) {
		return nil, pg.ErrMissingParameter
	}

	// validation parameters
	err := validationParams(params)
	if err != nil {
		return nil, err
	}

	// create uri and check credentials
	uri, ok := setUriChargeAndCheckCredentials(opts)
	if !ok {
		return nil, pg.ErrInvalidCredentials
	}

	// switch statement for handling queryParam or bodyParam
	var payload []byte
	switch params.(type) {
	case *CardToken:
		ct := params.(*CardToken)

		u := url.Values{}
		u.Set(clientKey, opts.ClientId)
		u.Set(cardCvv, ct.CardCvv)

		// use tokenId instead if exists
		if !reflect.ValueOf(ct.TokenID).IsZero() {
			u.Set(tokenId, ct.TokenID)
		} else {
			u.Set(cardNumber, ct.CardNumber)
			u.Set(cardExpMonth, ct.CardExpMonth)
			u.Set(cardExpYear, ct.CardExpYear)
		}

		payload = []byte(u.Encode())
	case *CardRegister:
		cr := params.(*CardRegister)

		u := url.Values{}
		u.Set(clientKey, opts.ClientId)
		u.Set(cardCvv, cr.CardCvv)
		u.Set(cardNumber, cr.CardNumber)
		u.Set(cardExpMonth, cr.CardExpMonth)
		u.Set(cardExpYear, cr.CardExpYear)

		payload = []byte(u.Encode())
	default:
		payload, _ = json.Marshal(params)
	}

	// create instance midtrans
	m := &midtrans{
		uri:    uri,
		opts:   opts,
		params: payload,
	}

	return m, nil
}

func setUriChargeAndCheckCredentials(opts *pg.Options) (string, bool) {
	var uri = sandBoxUri
	var env = pg.SandBox
	if opts.Environment == pg.Production {
		uri = productionUri
		env = pg.Production
	}

	// check valid credentials for each environment
	if !checkCredentials(opts.ServerKey, env) {
		return "", false
	}

	return uri, true
}

// checkCredentials credentials validation environment
func checkCredentials(c string, env pg.EnvironmentType) bool {
	c = strings.ToLower(c) // make lower case
	sp := strings.Split(c, "-")

	// environment sandBox will start with SB
	if env == pg.SandBox {
		return sp[0] == "sb"
	}

	return sp[0] != "sb"
}
