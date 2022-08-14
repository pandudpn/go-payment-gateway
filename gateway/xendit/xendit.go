package xendit

import (
	"encoding/json"
	"fmt"
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

	// pathParams parameter on URL Path
	pathParams string

	// params request to charge payment
	params []byte
}

func createChargeXendit(params interface{}, opts *pg.Options) (*xendit, error) {
	// return an error when params not exists
	if params == nil || (reflect.ValueOf(params).Kind() == reflect.Ptr && reflect.ValueOf(params).IsNil()) {
		return nil, pg.ErrMissingParameter
	}

	// validation parameters
	if reflect.ValueOf(params).Kind() == reflect.Ptr || reflect.ValueOf(params).Kind() == reflect.Struct {
		err := ValidationParams(params)
		if err != nil {
			return nil, err
		}
	}

	// create uri and check credentials
	ok := setUriChargeAndCheckCredentials(opts)
	if !ok {
		return nil, pg.ErrInvalidCredentials
	}

	var payload []byte
	var pathParams string

	switch param := params.(type) {
	case *UpdateVirtualAccountParam:
		pathParams = fmt.Sprintf("/%s", param.ID)
		payload, _ = json.Marshal(param)
	default:
		payload, _ = json.Marshal(param)
	}

	// create instance xendit
	m := &xendit{
		opts:       opts,
		pathParams: pathParams,
		params:     payload,
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
