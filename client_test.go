package pg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// mockProvider is a mock implementation of the Provider interface for testing
type mockProvider struct {
	name          string
	chargeResp    *ChargeResponse
	chargeErr     error
	statusResp    *PaymentStatus
	statusErr     error
	cancelErr     error
	webhookValid  bool
	webhookEvent  *WebhookEvent
	webhookErr    error
	webhookCalled bool
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) CreateCharge(ctx context.Context, params ChargeParams) (*ChargeResponse, error) {
	m.webhookCalled = true
	return m.chargeResp, m.chargeErr
}

func (m *mockProvider) GetStatus(ctx context.Context, orderID string) (*PaymentStatus, error) {
	return m.statusResp, m.statusErr
}

func (m *mockProvider) Cancel(ctx context.Context, orderID string) error {
	return m.cancelErr
}

func (m *mockProvider) VerifyWebhook(r *http.Request) bool {
	return m.webhookValid
}

func (m *mockProvider) ParseWebhook(r *http.Request) (*WebhookEvent, error) {
	return m.webhookEvent, m.webhookErr
}

// mockProviderFactory creates a mock provider for testing
func mockProviderFactory(cfg *ProviderConfig) (Provider, error) {
	return &mockProvider{name: "mock"}, nil
}

func TestRegisterProvider(t *testing.T) {
	// Register a mock provider
	RegisterProvider("mock-test", mockProviderFactory)

	providersMutex.RLock()
	_, ok := providers["mock-test"]
	providersMutex.RUnlock()

	if !ok {
		t.Error("provider was not registered")
	}

	// Clean up
	providersMutex.Lock()
	delete(providers, "mock-test")
	providersMutex.Unlock()
}

func TestNewClient(t *testing.T) {
	// Register mock provider
	RegisterProvider("test-provider", mockProviderFactory)
	defer func() {
		providersMutex.Lock()
		delete(providers, "test-provider")
		providersMutex.Unlock()
	}()

	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with options",
			opts: []Option{
				WithProvider("test-provider"),
				WithServerKey("test-key"),
				WithEnvironment("sandbox"),
			},
			wantErr: false,
		},
		{
			name: "missing provider",
			opts: []Option{
				WithProvider(""),
				WithServerKey("test-key"),
			},
			wantErr: true,
			errMsg:  "invalid configuration",
		},
		{
			name: "missing server key",
			opts: []Option{
				WithProvider("test-provider"),
			},
			wantErr: true,
			errMsg:  "invalid configuration",
		},
		{
			name: "unsupported provider",
			opts: []Option{
				WithProvider("unsupported"),
				WithServerKey("test-key"),
			},
			wantErr: true,
			errMsg:  "unsupported provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error message = %v, want contain %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if client == nil {
					t.Error("expected client but got nil")
				}
			}
		})
	}
}

func TestNewClientFromEnv(t *testing.T) {
	// Register mock provider
	RegisterProvider("env-provider", mockProviderFactory)
	defer func() {
		providersMutex.Lock()
		delete(providers, "env-provider")
		providersMutex.Unlock()
	}()

	// Set environment variables
	os.Setenv(EnvProvider, "env-provider")
	os.Setenv(EnvServerKey, "env-key")
	os.Setenv(EnvEnv, "sandbox")
	defer func() {
		os.Unsetenv(EnvProvider)
		os.Unsetenv(EnvServerKey)
		os.Unsetenv(EnvEnv)
	}()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.GetProvider() != "env-provider" {
		t.Errorf("provider = %v, want env-provider", client.GetProvider())
	}

	if client.GetConfig().ServerKey != "env-key" {
		t.Errorf("server key = %v, want env-key", client.GetConfig().ServerKey)
	}
}

func TestClient_CreateCharge(t *testing.T) {
	mock := &mockProvider{
		name: "mock",
		chargeResp: &ChargeResponse{
			TransactionID: "txn-123",
			OrderID:       "ORDER-001",
			Amount:        50000,
			Status:        StatusPending,
			PaymentURL:    "https://example.com/pay",
		},
	}

	client := &Client{
		provider: mock,
		config:   &Config{},
	}

	resp, err := client.CreateCharge(context.Background(), ChargeParams{
		OrderID: "ORDER-001",
		Amount:  50000,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.TransactionID != "txn-123" {
		t.Errorf("TransactionID = %v, want txn-123", resp.TransactionID)
	}

	if !mock.webhookCalled {
		t.Error("provider was not called")
	}
}

func TestClient_GetStatus(t *testing.T) {
	mock := &mockProvider{
		name: "mock",
		statusResp: &PaymentStatus{
			TransactionID: "txn-123",
			OrderID:       "ORDER-001",
			Status:        StatusSuccess,
			Amount:        50000,
			PaidAmount:    50000,
		},
	}

	client := &Client{
		provider: mock,
		config:   &Config{},
	}

	resp, err := client.GetStatus(context.Background(), "ORDER-001")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.Status != StatusSuccess {
		t.Errorf("Status = %v, want SUCCESS", resp.Status)
	}
}

func TestClient_Cancel(t *testing.T) {
	tests := []struct {
		name      string
		cancelErr error
		wantErr   bool
	}{
		{
			name:      "successful cancel",
			cancelErr: nil,
			wantErr:   false,
		},
		{
			name:      "cancel failed",
			cancelErr: ErrTransactionFailed,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{
				name:      "mock",
				cancelErr: tt.cancelErr,
			}

			client := &Client{
				provider: mock,
				config:   &Config{},
			}

			err := client.Cancel(context.Background(), "ORDER-001")
			if (err != nil) != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ParseWebhook(t *testing.T) {
	tests := []struct {
		name         string
		webhookValid bool
		webhookEvent *WebhookEvent
		webhookErr   error
		wantErr      bool
	}{
		{
			name:         "valid webhook",
			webhookValid: true,
			webhookEvent: &WebhookEvent{
				OrderID:       "ORDER-001",
				TransactionID: "txn-123",
				Status:        StatusSuccess,
				Amount:        50000,
				EventType:     EventPaymentCompleted,
			},
			wantErr: false,
		},
		{
			name:         "invalid signature",
			webhookValid: false,
			wantErr:      true,
		},
		{
			name:         "parse error",
			webhookValid: true,
			webhookErr:   ErrInvalidPayload,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{
				name:         "mock",
				webhookValid: tt.webhookValid,
				webhookEvent: tt.webhookEvent,
				webhookErr:   tt.webhookErr,
			}

			client := &Client{
				provider: mock,
				config:   &Config{},
			}

			req := httptest.NewRequest("POST", "/webhook", nil)

			event, err := client.ParseWebhook(req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if event == nil && tt.webhookEvent != nil {
					t.Error("expected event but got nil")
				}
			}
		})
	}
}

func TestClient_GetProvider(t *testing.T) {
	client := &Client{
		config: &Config{
			Provider: "midtrans",
		},
	}

	if client.GetProvider() != "midtrans" {
		t.Errorf("GetProvider() = %v, want midtrans", client.GetProvider())
	}
}

func TestClient_GetConfig(t *testing.T) {
	cfg := &Config{
		Provider:    "xendit",
		Environment: Production,
		ServerKey:   "test-key",
	}

	client := &Client{
		config: cfg,
	}

	result := client.GetConfig()
	if result.Provider != cfg.Provider {
		t.Errorf("GetConfig().Provider = %v, want %v", result.Provider, cfg.Provider)
	}
}

func TestClient_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentType
		expected bool
	}{
		{
			name:     "production",
			env:      Production,
			expected: true,
		},
		{
			name:     "sandbox",
			env:      SandBox,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				config: &Config{
					Environment: tt.env,
				},
			}

			if got := client.IsProduction(); got != tt.expected {
				t.Errorf("IsProduction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClient_IsSandbox(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentType
		expected bool
	}{
		{
			name:     "production",
			env:      Production,
			expected: false,
		},
		{
			name:     "sandbox",
			env:      SandBox,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				config: &Config{
					Environment: tt.env,
				},
			}

			if got := client.IsSandbox(); got != tt.expected {
				t.Errorf("IsSandbox() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCreateProvider(t *testing.T) {
	// Register mock provider
	RegisterProvider("create-test", mockProviderFactory)
	defer func() {
		providersMutex.Lock()
		delete(providers, "create-test")
		providersMutex.Unlock()
	}()

	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid provider",
			cfg: &Config{
				Provider:  "create-test",
				ServerKey: "test-key",
			},
			wantErr: false,
		},
		{
			name: "unsupported provider",
			cfg: &Config{
				Provider:  "unsupported",
				ServerKey: "test-key",
			},
			wantErr: true,
			errMsg:  "unsupported provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := createProvider(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error message = %v, want contain %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if provider == nil {
					t.Error("expected provider but got nil")
				}
			}
		})
	}
}
