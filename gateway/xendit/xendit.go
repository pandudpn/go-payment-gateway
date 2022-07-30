package xendit

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/pandudpn/go-payment-gateway"
)

const (
	// uri api
	uri = "https://api.xendit.co"

	// uri for charge transaction with payment method e-wallet
	ewalletUri = "/ewallets/charges"

	// uri for create virtual account
	vaUri = "/callback_virtual_accounts"
)

type xendit struct {
	// opts configuration or credentials for http api call
	opts *pg.Options

	// uri http target to charge payments
	uri string

	// params request to charge payment
	params []byte
}

func createChargeXendit(params interface{}, opts *pg.Options) (*xendit, error) {
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
	ok := setUriChargeAndCheckCredentials(opts)
	if !ok {
		return nil, pg.ErrInvalidCredentials
	}

	payload, _ := json.Marshal(params)

	// create instance xendit
	m := &xendit{
		opts:   opts,
		params: payload,
	}

	return m, nil
}

func setUriChargeAndCheckCredentials(opts *pg.Options) bool {
	var env = pg.SandBox
	if opts.Environment == pg.Production {
		env = pg.Production
	}

	// check valid credentials for each environment
	if !checkCredentials(opts.ServerKey, env) {
		return false
	}

	return true
}

// checkCredentials credentials validation environment
func checkCredentials(c string, env pg.EnvironmentType) bool {
	c = strings.ToLower(c) // make lower case

	// environment sandBox will start with SB
	if env == pg.SandBox {
		return strings.Contains(c, "xnd_development_")
	}

	return strings.Contains(c, "xnd_production_")
}
