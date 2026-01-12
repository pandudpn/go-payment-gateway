package pg

import (
"os"
"testing"
"time"
)

func TestWithProvider(t *testing.T) {
	cfg := &Config{}
	WithProvider("xendit")(cfg)

	if cfg.Provider != "xendit" {
		t.Errorf("Provider = %v, want xendit", cfg.Provider)
	}
}

func TestWithEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected EnvironmentType
	}{
		{"production", "production", Production},
		{"Production", "Production", SandBox},
		{"sandbox", "sandbox", SandBox},
		{"invalid", "invalid", SandBox},
		{"", "", SandBox},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			WithEnvironment(tt.input)(cfg)

			if cfg.Environment != tt.expected {
				t.Errorf("Environment = %v, want %v", cfg.Environment, tt.expected)
			}
		})
	}
}

func TestWithServerKey(t *testing.T) {
	cfg := &Config{}
	key := "test-server-key"
	WithServerKey(key)(cfg)

	if cfg.ServerKey != key {
		t.Errorf("ServerKey = %v, want %v", cfg.ServerKey, key)
	}
}

func TestWithClientKey(t *testing.T) {
	cfg := &Config{}
	key := "test-client-key"
	WithClientKey(key)(cfg)

	if cfg.ClientKey != key {
		t.Errorf("ClientKey = %v, want %v", cfg.ClientKey, key)
	}
}

func TestWithMerchantID(t *testing.T) {
	cfg := &Config{}
	id := "merchant-123"
	WithMerchantID(id)(cfg)

	if cfg.MerchantID != id {
		t.Errorf("MerchantID = %v, want %v", cfg.MerchantID, id)
	}
}

func TestWithTimeout(t *testing.T) {
	cfg := &Config{}
	timeout := 60 * time.Second
	WithTimeout(timeout)(cfg)

	if cfg.Timeout != timeout {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, timeout)
	}
}

func TestWithSnap(t *testing.T) {
	cfg := &Config{}
	WithSnap()(cfg)

	if !cfg.SnapMode {
		t.Error("SnapMode should be true")
	}
}

func TestWithLogging(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
	}{
		{"enabled", true},
		{"disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			WithLogging(tt.enabled)(cfg)

			if cfg.LogEnabled != tt.enabled {
				t.Errorf("LogEnabled = %v, want %v", cfg.LogEnabled, tt.enabled)
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv(EnvProvider, "midtrans")
	os.Setenv(EnvServerKey, "server-key-123")
	os.Setenv(EnvClientKey, "client-key-123")
	os.Setenv(EnvMerchantID, "merchant-123")
	os.Setenv(EnvEnv, "production")
	os.Setenv(EnvSnap, "true")
	os.Setenv(EnvLogging, "false")
	defer func() {
		os.Unsetenv(EnvProvider)
		os.Unsetenv(EnvServerKey)
		os.Unsetenv(EnvClientKey)
		os.Unsetenv(EnvMerchantID)
		os.Unsetenv(EnvEnv)
		os.Unsetenv(EnvSnap)
		os.Unsetenv(EnvLogging)
	}()

	cfg := LoadConfigFromEnv()

	if cfg.Provider != "midtrans" {
		t.Errorf("Provider = %v, want midtrans", cfg.Provider)
	}

	if cfg.ServerKey != "server-key-123" {
		t.Errorf("ServerKey = %v, want server-key-123", cfg.ServerKey)
	}

	if cfg.ClientKey != "client-key-123" {
		t.Errorf("ClientKey = %v, want client-key-123", cfg.ClientKey)
	}

	if cfg.MerchantID != "merchant-123" {
		t.Errorf("MerchantID = %v, want merchant-123", cfg.MerchantID)
	}

	if cfg.Environment != Production {
		t.Errorf("Environment = %v, want production", cfg.Environment)
	}

	if !cfg.SnapMode {
		t.Error("SnapMode should be true")
	}

	if cfg.LogEnabled {
		t.Error("LogEnabled should be false")
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", cfg.Timeout)
	}
}

func TestLoadConfigFromEnvDefaults(t *testing.T) {
	// Ensure no env vars are set
	os.Unsetenv(EnvProvider)
	os.Unsetenv(EnvServerKey)
	os.Unsetenv(EnvClientKey)
	os.Unsetenv(EnvMerchantID)
	os.Unsetenv(EnvEnv)
	os.Unsetenv(EnvSnap)
	os.Unsetenv(EnvLogging)

	cfg := LoadConfigFromEnv()

	if cfg.Provider != "" {
		t.Errorf("Provider should be empty, got %v", cfg.Provider)
	}

	if cfg.Environment != SandBox {
		t.Errorf("Environment = %v, want SandBox", cfg.Environment)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", cfg.Timeout)
	}

	if cfg.SnapMode {
		t.Error("SnapMode should be false by default")
	}

	if !cfg.LogEnabled {
		t.Error("LogEnabled should be true by default")
	}
}

func TestApplyOptions(t *testing.T) {
	cfg := &Config{}

	result := ApplyOptions(cfg,
		WithProvider("doku"),
		WithServerKey("test-key"),
		WithSnap(),
	)

	if result.Provider != "doku" {
		t.Errorf("Provider = %v, want doku", result.Provider)
	}

	if result.ServerKey != "test-key" {
		t.Errorf("ServerKey = %v, want test-key", result.ServerKey)
	}

	if !result.SnapMode {
		t.Error("SnapMode should be true")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Provider != ProviderMidtrans {
		t.Errorf("Provider = %v, want %v", cfg.Provider, ProviderMidtrans)
	}

	if cfg.Environment != SandBox {
		t.Errorf("Environment = %v, want SandBox", cfg.Environment)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", cfg.Timeout)
	}

	if cfg.SnapMode {
		t.Error("SnapMode should be false by default")
	}

	if !cfg.LogEnabled {
		t.Error("LogEnabled should be true by default")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Provider:  "midtrans",
				ServerKey: "test-key",
			},
			wantErr: false,
		},
		{
			name: "missing provider",
			cfg: &Config{
				ServerKey: "test-key",
			},
			wantErr: true,
		},
		{
			name: "missing server key",
			cfg: &Config{
				Provider: "midtrans",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			cfg:     &Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentType
		expected bool
	}{
		{"production", Production, true},
		{"sandbox", SandBox, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Environment: tt.env}
			if got := cfg.IsProduction(); got != tt.expected {
				t.Errorf("IsProduction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_IsSandbox(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentType
		expected bool
	}{
		{"production", Production, false},
		{"sandbox", SandBox, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Environment: tt.env}
			if got := cfg.IsSandbox(); got != tt.expected {
				t.Errorf("IsSandbox() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_IsSnap(t *testing.T) {
	tests := []struct {
		name     string
		snapMode bool
		expected bool
	}{
		{"snap enabled", true, true},
		{"snap disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{SnapMode: tt.snapMode}
			if got := cfg.IsSnap(); got != tt.expected {
				t.Errorf("IsSnap() = %v, want %v", got, tt.expected)
			}
		})
	}
}
