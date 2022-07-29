package midtrans

import (
	"encoding/json"
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
	if params == nil {
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

	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
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

	return uri + chargeUri, true
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
