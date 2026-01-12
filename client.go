package pg

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Version is the library version
const Version = "v1.0.0"

// ProviderFactory creates a provider from configuration
type ProviderFactory func(cfg *ProviderConfig) (Provider, error)

// Provider is the interface for payment gateway providers
type Provider interface {
	Name() string
	CreateCharge(ctx context.Context, params ChargeParams) (*ChargeResponse, error)
	GetStatus(ctx context.Context, orderID string) (*PaymentStatus, error)
	Cancel(ctx context.Context, orderID string) error
	VerifyWebhook(r *http.Request) bool
	ParseWebhook(r *http.Request) (*WebhookEvent, error)
}

// ProviderConfig holds the configuration for a provider
type ProviderConfig struct {
	Environment string
	ServerKey   string
	ClientKey   string
	MerchantID  string
	Timeout     int
	SnapMode    bool
	LogEnabled  bool
}

var (
	providers      = make(map[string]ProviderFactory)
	providersMutex sync.RWMutex
)

// RegisterProvider registers a provider factory
// This should be called from the provider's init() function
func RegisterProvider(name string, factory ProviderFactory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	providers[name] = factory
}

// Client is the main entry point for payment operations
type Client struct {
	provider Provider
	config   *Config
}

// NewClient creates a new payment client with the given options
// If no options are provided, it will load configuration from environment variables
func NewClient(opts ...Option) (*Client, error) {
	// Start with default config
	cfg := DefaultConfig()

	// Load from environment first (so env vars can be defaults)
	envCfg := LoadConfigFromEnv()
	if envCfg.Provider != "" {
		cfg.Provider = envCfg.Provider
	}
	if envCfg.ServerKey != "" {
		cfg.ServerKey = envCfg.ServerKey
	}
	if envCfg.ClientKey != "" {
		cfg.ClientKey = envCfg.ClientKey
	}
	if envCfg.MerchantID != "" {
		cfg.MerchantID = envCfg.MerchantID
	}

	// Apply explicit options (overrides env vars)
	cfg = ApplyOptions(cfg, opts...)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create provider based on config
	p, err := createProvider(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		provider: p,
		config:   cfg,
	}, nil
}

// createProvider creates a provider based on the configuration
func createProvider(cfg *Config) (Provider, error) {
	providersMutex.RLock()
	factory, ok := providers[cfg.Provider]
	providersMutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	providerCfg := &ProviderConfig{
		Environment: string(cfg.Environment),
		ServerKey:   cfg.ServerKey,
		ClientKey:   cfg.ClientKey,
		MerchantID:  cfg.MerchantID,
		Timeout:     int(cfg.Timeout.Seconds()),
		SnapMode:    cfg.SnapMode,
		LogEnabled:  cfg.LogEnabled,
	}

	return factory(providerCfg)
}

// CreateCharge creates a new payment transaction
func (c *Client) CreateCharge(ctx context.Context, params ChargeParams) (*ChargeResponse, error) {
	return c.provider.CreateCharge(ctx, params)
}

// GetStatus retrieves the status of a payment transaction
func (c *Client) GetStatus(ctx context.Context, orderID string) (*PaymentStatus, error) {
	return c.provider.GetStatus(ctx, orderID)
}

// Cancel cancels a payment transaction
func (c *Client) Cancel(ctx context.Context, orderID string) error {
	return c.provider.Cancel(ctx, orderID)
}

// ParseWebhook parses and verifies a webhook notification
func (c *Client) ParseWebhook(r *http.Request) (*WebhookEvent, error) {
	// Verify signature first
	if !c.provider.VerifyWebhook(r) {
		return nil, ErrInvalidSignature
	}

	// Parse webhook
	return c.provider.ParseWebhook(r)
}

// GetProvider returns the name of the current provider
func (c *Client) GetProvider() string {
	return c.config.Provider
}

// GetConfig returns a copy of the client configuration
func (c *Client) GetConfig() Config {
	return *c.config
}

// IsProduction returns true if the client is in production mode
func (c *Client) IsProduction() bool {
	return c.config.IsProduction()
}

// IsSandbox returns true if the client is in sandbox mode
func (c *Client) IsSandbox() bool {
	return c.config.IsSandbox()
}
