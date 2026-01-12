package pg

import (
	"os"
	"time"
)

// Config holds the configuration for the payment client
type Config struct {
	// Provider is the payment gateway provider name (midtrans, xendit, doku)
	Provider string

	// Environment is the environment (sandbox, production)
	Environment EnvironmentType

	// ServerKey is the server key for authentication
	ServerKey string

	// ClientKey is the client key (provider-specific)
	ClientKey string

	// MerchantID is the merchant ID (provider-specific)
	MerchantID string

	// Timeout is the HTTP request timeout
	Timeout time.Duration

	// SnapMode indicates if SNAP mode is enabled (for Midtrans)
	SnapMode bool

	// LogEnabled indicates if logging is enabled
	LogEnabled bool
}

// Option is a function that configures the client
type Option func(*Config)

// WithProvider sets the payment gateway provider
func WithProvider(provider string) Option {
	return func(c *Config) {
		c.Provider = provider
	}
}

// WithEnvironment sets the environment
func WithEnvironment(env string) Option {
	return func(c *Config) {
		if env == string(Production) {
			c.Environment = Production
		} else {
			c.Environment = SandBox
		}
	}
}

// WithServerKey sets the server key
func WithServerKey(key string) Option {
	return func(c *Config) {
		c.ServerKey = key
	}
}

// WithClientKey sets the client key
func WithClientKey(key string) Option {
	return func(c *Config) {
		c.ClientKey = key
	}
}

// WithMerchantID sets the merchant ID
func WithMerchantID(id string) Option {
	return func(c *Config) {
		c.MerchantID = id
	}
}

// WithTimeout sets the HTTP request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithSnap enables SNAP mode (for Midtrans)
func WithSnap() Option {
	return func(c *Config) {
		c.SnapMode = true
	}
}

// WithLogging enables or disables logging
func WithLogging(enabled bool) Option {
	return func(c *Config) {
		c.LogEnabled = enabled
	}
}

// Environment variable names
const (
	EnvProvider   = "PAYMENT_PROVIDER"
	EnvEnv        = "PAYMENT_ENV"
	EnvServerKey  = "PAYMENT_SERVER_KEY"
	EnvClientKey  = "PAYMENT_CLIENT_KEY"
	EnvMerchantID = "PAYMENT_MERCHANT_ID"
	EnvTimeout    = "PAYMENT_TIMEOUT"
	EnvSnap       = "PAYMENT_SNAP"
	EnvLogging    = "PAYMENT_LOGGING"
)

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	cfg := &Config{
		Provider:    os.Getenv(EnvProvider),
		Environment: SandBox,
		ServerKey:   os.Getenv(EnvServerKey),
		ClientKey:   os.Getenv(EnvClientKey),
		MerchantID:  os.Getenv(EnvMerchantID),
		Timeout:     30 * time.Second,
		SnapMode:    os.Getenv(EnvSnap) == "true",
		LogEnabled:  os.Getenv(EnvLogging) != "false",
	}

	if env := os.Getenv(EnvEnv); env == string(Production) {
		cfg.Environment = Production
	}

	return cfg
}

// ApplyOptions applies a list of options to a config
func ApplyOptions(cfg *Config, opts ...Option) *Config {
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Provider:    ProviderMidtrans,
		Environment: SandBox,
		Timeout:     30 * time.Second,
		SnapMode:    false,
		LogEnabled:  true,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Provider == "" {
		return NewRequiredFieldError("Provider")
	}
	if c.ServerKey == "" {
		return NewRequiredFieldError("ServerKey")
	}
	return nil
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == Production
}

// IsSandbox returns true if the environment is sandbox
func (c *Config) IsSandbox() bool {
	return c.Environment == SandBox
}

// IsSnap returns true if SNAP mode is enabled
func (c *Config) IsSnap() bool {
	return c.SnapMode
}
