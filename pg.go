package pg

import (
	"fmt"

	"github.com/pandudpn/go-payment-gateway/utils"
)

const (
	Midtrans string = "midtrans"
	Xendit   string = "xendit"
	Oy       string = "oy! indonesia"
)

type PG struct {
	// configuration payment gateway
	config *Config

	// midtrans configuration
	midtrans *midtrans

	// xendit configuration
	xendit *xendit
}

// New creates a new Payment Gateway instance
func New(cfg ...*Config) *PG {
	p := &PG{
		config:   DefaultConfig,
		midtrans: &midtrans{uri: mdUriSandbox, credentials: &Credentials{}},
		xendit:   &xendit{uri: xndUri, credentials: &Credentials{}},
	}

	if len(cfg) > 0 {
		p.config = cfg[0]
	}
	// disable logging
	if !p.config.Logging {
		utils.DisableLogging()
	}

	// when configuration from parameter and env is production
	// set all configuration to production
	if p.config.Environment == Production {
		p.SetProduction()
	}

	return p
}

// SetProduction environment and base url gateway
func (p *PG) SetProduction() {
	// re-value env with Production
	p.config.Environment = Production

	// set production uri to all gateway
	p.midtrans.uri = mdUriProduction
}

// SetMidtransCredentials for access Midtrans Core API
func (p *PG) SetMidtransCredentials(clientSecret string) error {
	// check valid client secret
	if !checkCredentials(p.config.Environment, clientSecret, Midtrans) {
		return fmt.Errorf("%s credentials invalid for environment %s", Midtrans, p.config.Environment)
	}

	p.midtrans.credentials.ClientSecret = clientSecret
	return nil
}

// SetXenditCredentials for access Xendit API
func (p *PG) SetXenditCredentials(clientSecret string) error {
	// check valid client secret
	if !checkCredentials(p.config.Environment, clientSecret, Xendit) {
		return fmt.Errorf("%s credentials invalid for environment %s", Xendit, p.config.Environment)
	}

	p.xendit.credentials.ClientSecret = clientSecret
	return nil
}

// Midtrans get instance of midtrans gateway
func (p *PG) Midtrans() (*midtrans, error) {
	if p.midtrans == nil {
		utils.Log.Errorf("credentials midtrans is nil")
		return nil, ErrNilCredentials
	}

	return p.midtrans, nil
}

// Xendit get instance of xendit gateway
func (p *PG) Xendit() (*xendit, error) {
	if p.xendit == nil {
		utils.Log.Errorf("credentials xendit is nil")
		return nil, ErrNilCredentials
	}

	return p.xendit, nil
}
