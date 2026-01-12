# Unified Payment Gateway - Specification Document

**Version**: 2.0.0
**Date**: 2025-01-12
**Status**: Draft

---

## 1. Overview

### 1.1 Purpose
Centralized payment gateway library for Indonesian payment providers with unified interface. Users interact with a single API while the library handles provider-specific implementations internally.

### 1.2 Supported Providers (Phase 1)
| Provider | Status | Notes |
|----------|--------|-------|
| Midtrans | ✅ Phase 1 | With SNAP support |
| Xendit | ✅ Phase 1 | Full support |
| Doku | ✅ Phase 1 | Full support |
| Duitku | ⏳ Phase 2 | Future |
| Faspay | ⏳ Phase 2 | Future |
| Espay | ⏳ Phase 2 | Future |

### 1.3 Supported Payment Types
| Category | Types |
|----------|-------|
| E-Wallet | GoPay, OVO, DANA, ShopeePay, LinkAja |
| Virtual Account | BCA, BNI, BRI, Mandiri, Permata, CIMB |
| QRIS | QRIS (all providers) |
| Credit Card | Visa, Mastercard, JCB |
| Retail | Alfamart, Indomaret |

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    User Application Layer                   │
│  import "github.com/pandudpn/go-payment-gateway/pg"        │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      Public API Layer (pg)                  │
│  • Client  • Types  • Options  • Constants  • Errors       │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                   Internal Provider Layer                   │
│  internal/provider/provider.go (Provider Interface)         │
└────────────────────────────┬────────────────────────────────┘
                             │
                ┌────────────┼────────────┐
                ▼            ▼            ▼
          ┌─────────┐  ┌─────────┐  ┌─────────┐
          │Midtrans │  │ Xendit  │  │  Doku   │
          │ Provider│  │Provider │  │Provider │
          └─────────┘  └─────────┘  └─────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      Internal Utils                         │
│  • Validator  • Signature  • Mapper  • Config               │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Public API Specification

### 3.1 Client Initialization

```go
package pg

// Client is the main entry point for payment operations
type Client struct {
    provider Provider
    config   *Config
}

// Option pattern for client configuration
type Option func(*Config)

// NewClient creates a new payment client
func NewClient(opts ...Option) (*Client, error)

// Functional Options
func WithProvider(name string) Option
func WithSnap() Option
func WithEnvironment(env string) Option
func WithServerKey(key string) Option
func WithClientKey(key string) Option
func WithMerchantID(id string) Option
func WithTimeout(timeout time.Duration) Option
```

**Usage Example:**
```go
// Explicit configuration
client := pg.NewClient(
    pg.WithProvider("midtrans"),
    pg.WithSnap(),
    pg.WithEnvironment("sandbox"),
    pg.WithServerKey(os.Getenv("PAYMENT_SERVER_KEY")),
    pg.WithClientKey(os.Getenv("PAYMENT_CLIENT_KEY")),
    pg.WithMerchantID(os.Getenv("PAYMENT_MERCHANT_ID")),
    pg.WithTimeout(30 * time.Second),
)

// Or auto-load from environment variables
client := pg.NewClient() // Reads PAYMENT_PROVIDER, PAYMENT_ENV, etc.
```

### 3.2 Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| PAYMENT_PROVIDER | Yes | - | Provider name: midtrans, xendit, doku |
| PAYMENT_ENV | No | sandbox | Environment: sandbox, production |
| PAYMENT_SERVER_KEY | Yes | - | Server key for authentication |
| PAYMENT_CLIENT_KEY | No | - | Client key (provider-specific) |
| PAYMENT_MERCHANT_ID | No | - | Merchant ID (provider-specific) |
| PAYMENT_TIMEOUT | No | 30s | Request timeout in seconds |

---

## 4. Data Types Specification

### 4.1 Request Types

```go
// ChargeParams - Unified charge request parameters
type ChargeParams struct {
    // Required fields
    OrderID     string      `json:"order_id"`
    Amount      int64       `json:"amount"`
    PaymentType PaymentType `json:"payment_type"`

    // Customer information
    Customer    Customer    `json:"customer"`

    // Optional fields
    Items       []Item      `json:"items,omitempty"`
    Description string      `json:"description,omitempty"`
    ExpiryTime  time.Time   `json:"expiry_time,omitempty"`

    // Callback URLs
    CallbackURL string      `json:"callback_url,omitempty"`
    ReturnURL   string      `json:"return_url,omitempty"`

    // Provider-specific parameters (not mapped to unified fields)
    Custom      map[string]interface{} `json:"-"`
}

// Customer represents customer information
type Customer struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
}

// Item represents transaction item
type Item struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Price    int64   `json:"price"`
    Quantity int64   `json:"quantity"`
    Category string  `json:"category,omitempty"`
}
```

### 4.2 Response Types

```go
// ChargeResponse - Unified charge response
type ChargeResponse struct {
    // Unified fields
    TransactionID   string    `json:"transaction_id"`
    OrderID         string    `json:"order_id"`
    Amount          int64     `json:"amount"`
    Status          Status    `json:"status"`
    PaymentURL      string    `json:"payment_url,omitempty"`
    QRString        string    `json:"qr_string,omitempty"`
    ExpiryTime      time.Time `json:"expiry_time,omitempty"`
    CreatedAt       time.Time `json:"created_at"`

    // Raw response from provider (for provider-specific fields)
    Raw             map[string]interface{} `json:"-"`
}

// PaymentStatus - Unified payment status response
type PaymentStatus struct {
    TransactionID   string     `json:"transaction_id"`
    OrderID         string     `json:"order_id"`
    Status          Status     `json:"status"`
    Amount          int64      `json:"amount"`
    PaidAmount      int64      `json:"paid_amount"`
    PaymentType     PaymentType `json:"payment_type"`
    PaidAt          *time.Time `json:"paid_at,omitempty"`

    // Raw response
    Raw             map[string]interface{} `json:"-"`
}
```

### 4.3 Webhook Types

```go
// WebhookEvent - Unified webhook event
type WebhookEvent struct {
    // Unified fields
    OrderID       string      `json:"order_id"`
    TransactionID string      `json:"transaction_id"`
    Status        Status      `json:"status"`
    Amount        int64       `json:"amount"`
    PaymentType   PaymentType `json:"payment_type"`
    EventType     string      `json:"event_type"`
    Timestamp     time.Time   `json:"timestamp"`

    // Raw from provider
    Raw           map[string]interface{} `json:"-"`
}

// Webhook notification
func (c *Client) ParseWebhook(r *http.Request) (*WebhookEvent, error)
```

---

## 5. Enumerations

### 5.1 Payment Types

```go
type PaymentType string

const (
    // E-Wallet
    PaymentTypeGoPay     PaymentType = "GOPAY"
    PaymentTypeOVO       PaymentType = "OVO"
    PaymentTypeDANA      PaymentType = "DANA"
    PaymentTypeShopeePay PaymentType = "SHOPEEPAY"
    PaymentTypeLinkAja   PaymentType = "LINKAJA"

    // Virtual Account
    PaymentTypeVABCA      PaymentType = "VA_BCA"
    PaymentTypeVABNI      PaymentType = "VA_BNI"
    PaymentTypeVABRI      PaymentType = "VA_BRI"
    PaymentTypeVAMandiri  PaymentType = "VA_MANDIRI"
    PaymentTypeVAPermata  PaymentType = "VA_PERMATA"
    PaymentTypeVACIMB     PaymentType = "VA_CIMB"

    // QRIS
    PaymentTypeQRIS       PaymentType = "QRIS"

    // Credit Card
    PaymentTypeCC         PaymentType = "CREDIT_CARD"

    // Retail
    PaymentTypeAlfamart   PaymentType = "ALFAMART"
    PaymentTypeIndomaret  PaymentType = "INDOMARET"
)
```

### 5.2 Transaction Status

```go
type Status string

const (
    StatusPending    Status = "PENDING"
    StatusProcessing Status = "PROCESSING"
    StatusSuccess    Status = "SUCCESS"
    StatusFailed     Status = "FAILED"
    StatusCancelled  Status = "CANCELLED"
    StatusExpired    Status = "EXPIRED"
)
```

### 5.3 Event Types

```go
const (
    EventPaymentCompleted   = "payment.completed"
    EventPaymentFailed      = "payment.failed"
    EventPaymentPending     = "payment.pending"
    EventPaymentExpired     = "payment.expired"
    EventPaymentCancelled   = "payment.cancelled"
)
```

---

## 6. Client Methods

### 6.1 Create Charge

```go
func (c *Client) CreateCharge(ctx context.Context, params ChargeParams) (*ChargeResponse, error)
```

**Error Responses:**
| Error | Condition |
|-------|-----------|
| `ErrMissingParameter` | Required field is empty |
| `ErrInvalidParameter` | Parameter validation failed |
| `ErrInvalidCredentials` | Credential validation failed |
| `ErrMinAmount` | Amount below minimum |
| `ErrDuplicateTransaction` | Order ID already exists |

### 6.2 Get Payment Status

```go
func (c *Client) GetStatus(ctx context.Context, orderID string) (*PaymentStatus, error)
```

### 6.3 Cancel Transaction

```go
func (c *Client) Cancel(ctx context.Context, orderID string) error
```

### 6.4 Parse Webhook

```go
func (c *Client) ParseWebhook(r *http.Request) (*WebhookEvent, error)
```

**Webhook Signature Validation:**
- Midtrans: SHA512(OrderID + Status + ServerKey)
- Xendit: Callback Token header validation
- Doku: Signature from config

**Error Responses:**
| Error | Condition |
|-------|-----------|
| `ErrInvalidSignature` | Webhook signature verification failed |
| `ErrInvalidPayload` | Payload parsing failed |

---

## 7. Error Handling

### 7.1 Unified Error Types

```go
var (
    // Base errors
    ErrMissingParameter   = errors.New("missing required parameter")
    ErrInvalidParameter   = errors.New("invalid parameter")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidSignature   = errors.New("invalid webhook signature")
    ErrInvalidPayload     = errors.New("invalid webhook payload")

    // Payment errors
    ErrMinAmount          = errors.New("minimum transaction amount is Rp10.000")
    ErrDuplicateTransaction = errors.New("duplicate transaction ID")
    ErrTransactionNotFound = errors.New("transaction not found")
    ErrTransactionFailed  = errors.New("transaction failed")

    // Network errors
    ErrTimeout            = errors.New("request timeout")
    ErrRateLimit          = errors.New("rate limit exceeded")
    ErrServiceUnavailable = errors.New("service unavailable")
)

// ProviderError wraps provider-specific errors
type ProviderError struct {
    Code    string
    Message string
    Provider string
    Raw     map[string]interface{}
}
```

### 7.2 Error Wrapping Pattern

```go
// Errors include context for debugging
charge, err := client.CreateCharge(ctx, params)
if err != nil {
    // Error contains:
    // - Unified error type
    // - Provider-specific code (if available)
    // - Raw provider response
    return fmt.Errorf("failed to create charge: %w", err)
}
```

---

## 8. Provider Interface (Internal)

```go
// internal/provider/provider.go

package provider

import (
    "context"
    "github.com/pandudpn/go-payment-gateway/pg"
)

// Provider defines the interface for payment gateway providers
type Provider interface {
    // Name returns the provider name
    Name() string

    // CreateCharge creates a new payment transaction
    CreateCharge(ctx context.Context, params pg.ChargeParams) (*pg.ChargeResponse, error)

    // GetStatus retrieves payment status
    GetStatus(ctx context.Context, orderID string) (*pg.PaymentStatus, error)

    // Cancel cancels a transaction
    Cancel(ctx context.Context, orderID string) error

    // VerifyWebhook verifies webhook signature
    VerifyWebhook(r *http.Request) bool

    // ParseWebhook parses webhook payload
    ParseWebhook(r *http.Request) (*pg.WebhookEvent, error)
}

// Config holds provider configuration
type Config struct {
    Environment string
    ServerKey   string
    ClientKey   string
    MerchantID  string
    Timeout     time.Duration
    SnapMode    bool
    IsSnap      bool // Alias for SnapMode
}
```

---

## 9. Mapping Specification

### 9.1 Unified → Provider Mapping

| Unified Field | Midtrans | Xendit | Doku |
|---------------|----------|--------|------|
| OrderID | transaction_details.order_id | external_id | transaction_data.order_id |
| Amount | transaction_details.gross_amount | amount | transaction_data.amount |
| Customer.ID | customer.id | customer.reference_id | customer.customer_id |
| Customer.Name | customer.first_name + last_name | customer.given_names | customer.customer_name |
| Customer.Email | customer.email | customer.email | customer.customer_email |
| Customer.Phone | customer.phone | customer.mobile_number | customer.customer_phone |
| PaymentType | payment_type | payment_method | payment.payment_method |

### 9.2 Provider → Unified Status Mapping

| Midtrans | Xendit | Doku | Unified |
|----------|--------|------|---------|
| pending | PENDING | REQUESTED | PENDING |
| processing | PROCESSING | PROCESSING | PROCESSING |
| settlement | SUCCEEDED | COMPLETED | SUCCESS |
| deny | FAILED | FAILED | FAILED |
| expire | EXPIRED | EXPIRED | EXPIRED |
| cancel | VOIDED | CANCELLED | CANCELLED |

---

## 10. Project Structure

```
go-payment-gateway/
├── pg.go                    # Main client
├── types.go                 # Public types
├── constants.go             # Enums
├── errors.go                # Errors
├── options.go               # Functional options
│
├── internal/                # Private - cannot be imported externally
│   ├── provider/
│   │   ├── provider.go      # Provider interface
│   │   ├── midtrans/
│   │   │   ├── midtrans.go  # Implementation
│   │   │   ├── mapper.go    # Mapping logic
│   │   │   ├── webhook.go   # Webhook handling
│   │   │   ├── types.go     # Internal types
│   │   │   └── constants.go # Provider constants
│   │   ├── xendit/
│   │   │   ├── xendit.go
│   │   │   ├── mapper.go
│   │   │   ├── webhook.go
│   │   │   ├── types.go
│   │   │   └── constants.go
│   │   └── doku/
│   │       ├── doku.go
│   │       ├── mapper.go
│   │       ├── webhook.go
│   │       ├── types.go
│   │       └── constants.go
│   │
│   ├── utils/
│   │   ├── validator.go     # Reusable validators
│   │   ├── signature.go     # Signature verification
│   │   └── mapper.go        # Generic mapper helpers
│   │
│   └── config/
│       ├── config.go        # Config struct
│       └── env.go           # Environment loading
│
├── example/
│   └── unified/
│       └── main.go          # Usage examples
│
└── go.mod
```

---

## 11. Implementation Phases

### Phase 1: Foundation (3-4 hours)
- [ ] Public API files (types, constants, errors, options)
- [ ] Internal provider interface
- [ ] Internal utils (validator, signature, mapper)
- [ ] Internal config (config, env loader)

### Phase 2: Midtrans (3-4 hours)
- [ ] Midtrans provider implementation
- [ ] Mapper (unified ↔ midtrans)
- [ ] Webhook handler with signature verification
- [ ] SNAP mode support
- [ ] Unit tests

### Phase 3: Xendit (3-4 hours)
- [ ] Xendit provider implementation
- [ ] Mapper (unified ↔ xendit)
- [ ] Webhook handler with callback token verification
- [ ] Unit tests

### Phase 4: Doku (4-5 hours)
- [ ] Doku provider implementation
- [ ] Mapper (unified ↔ doku)
- [ ] Webhook handler with signature verification
- [ ] Unit tests

### Phase 5: Polish (2-3 hours)
- [ ] Complete test coverage
- [ ] Documentation
- [ ] Usage examples
- [ ] README update

---

## 12. Testing Strategy

### 12.1 Unit Tests
- Validator functions
- Mapper functions (unified ↔ provider)
- Signature verification
- Config loading

### 12.2 Integration Tests
- Mock HTTP server per provider
- End-to-end charge flow
- Webhook parsing and verification
- Error handling scenarios

### 12.3 Test Coverage Target
- Minimum: 80%
- Target: 90%

---

## 13. Security Considerations

1. **Credential Storage**: Never log credentials
2. **Webhook Verification**: Always verify webhook signatures
3. **HTTPS Only**: Production API calls must use HTTPS
4. **Input Validation**: Validate all input parameters
5. **Secret Key Rotation**: Support for key rotation

---

## 14. Future Enhancements

- [ ] Phase 2 providers (Duitku, Faspay, Espay)
- [ ] Transaction history/listing
- [ ] Refund support
- [ ] Installment payments
- [ ] Recurring payments
- [ ] Payment method discovery (available methods per amount)
- [ ] Circuit breaker for provider failures
- [ ] Automatic fallback providers

---

**Document Version**: 1.0
**Last Updated**: 2025-01-12
