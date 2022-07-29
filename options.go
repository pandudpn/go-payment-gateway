package pg

import (
	"reflect"
)

var True bool = true

// Options is the wrap of the params needed for API Call
type Options struct {
	// Environment used for create transaction
	// Possible values are 'Production' or 'SandBox'
	//
	// Default: SandBox
	Environment EnvironmentType

	// ClientId is client_id from payment gateway
	// it can be apiKey or something like that
	ClientId string

	// ServerKey is client_secret from payment gateway
	ServerKey string

	// ApiCall interface to request outside
	//
	// Default: DefaultApiRequest
	ApiCall ApiRequestInterface

	// When set to true, this will log your request, response or error to stdout
	// Use logrus as logging
	//
	// Default: true
	Logging *bool
}

// DefaultOptions define all default value of configuration options payment gateway
var DefaultOptions = &Options{
	Environment: SandBox,
}

// NewOption will create an instance of configuration Options
func NewOption(opts ...*Options) (*Options, error) {
	opt := DefaultOptions

	if len(opts) > 0 {
		opt = opts[0]
	}

	// init standard logging
	NewLogger()

	// if environment not exists
	if reflect.ValueOf(opt.Environment).IsZero() {
		opt.Environment = SandBox
	}

	if opt.Logging == nil {
		opt.Logging = &True
	}

	// if logging is false or logging function is not exists
	// disable logging for handle an error
	if !*opt.Logging {
		DisableLogging()
	}

	// make sure clientId or clientSecret exists
	// if not exists, just given error ErrMissingCredentials
	if reflect.ValueOf(opt.ClientId).IsZero() || reflect.ValueOf(opt.ServerKey).IsZero() {
		return nil, ErrMissingCredentials
	}

	// make default apiRequestInterface if not exists
	if opt.ApiCall == nil {
		opt.ApiCall = DefaultApiRequest()
	}

	// return an instance of Options
	return opt, nil
}
