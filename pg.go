package pg

import (
	"fmt"
)

const (
	midtrans string = "midtrans"
	xendit   string = "xendit"
	oy       string = "oy! indonesia"
)

type PG struct {
	// configuration payment gateway
	config *Config
	
	// midtrans configuration
	midtrans *Midtrans
	
	// xendit configuration
	xendit *Xendit
}

// New creates a new Payment Gateway instance
func New(cfg ...*Config) *PG {
	p := &PG{
		config:   DefaultConfig,
		midtrans: &Midtrans{uri: mdUriSandbox, credentials: &Credentials{}},
		xendit:   &Xendit{uri: xndUri, credentials: &Credentials{}},
	}
	
	if len(cfg) > 0 {
		p.config = cfg[0]
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
	if !checkCredentials(p.config.Environment, clientSecret, midtrans) {
		return fmt.Errorf("%s credentials invalid for environment %s", midtrans, p.config.Environment)
	}
	
	p.midtrans.credentials.ClientSecret = clientSecret
	return nil
}

// SetXenditCredentials for access Xendit API
func (p *PG) SetXenditCredentials(clientSecret string) error {
	// check valid client secret
	if !checkCredentials(p.config.Environment, clientSecret, xendit) {
		return fmt.Errorf("%s credentials invalid for environment %s", xendit, p.config.Environment)
	}
	
	p.xendit.credentials.ClientSecret = clientSecret
	return nil
}
