package pg

import (
	"errors"
	"testing"
)

func TestFieldError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *FieldError
		expected string
	}{
		{
			name: "with message",
			err: &FieldError{
				Field:   "Email",
				Message: "must be valid",
			},
			expected: "Email: must be valid",
		},
		{
			name: "without message",
			err: &FieldError{
				Field: "Phone",
			},
			expected: "Phone is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	err := &FieldError{
		Field: "Test",
		Err:   baseErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestNewFieldError(t *testing.T) {
	err := NewFieldError("Email", "invalid format")

	if err.Field != "Email" {
		t.Errorf("Field = %v, want Email", err.Field)
	}

	if err.Message != "invalid format" {
		t.Errorf("Message = %v, want invalid format", err.Message)
	}

	if err.Err != ErrInvalidParameter {
		t.Errorf("Err = %v, want %v", err.Err, ErrInvalidParameter)
	}
}

func TestNewRequiredFieldError(t *testing.T) {
	err := NewRequiredFieldError("OrderID")

	if err.Field != "OrderID" {
		t.Errorf("Field = %v, want OrderID", err.Field)
	}

	if err.Message != "is required" {
		t.Errorf("Message = %v, want 'is required'", err.Message)
	}

	if err.Err != ErrMissingParameter {
		t.Errorf("Err = %v, want %v", err.Err, ErrMissingParameter)
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		errors   []*FieldError
		expected string
	}{
		{
			name:     "no errors",
			errors:   []*FieldError{},
			expected: "validation failed",
		},
		{
			name: "single error",
			errors: []*FieldError{
				{Field: "Email", Message: "invalid"},
			},
			expected: "Email: invalid",
		},
		{
			name: "multiple errors",
			errors: []*FieldError{
				{Field: "Email", Message: "invalid"},
				{Field: "Phone", Message: "invalid"},
			},
			expected: "validation failed: 2 errors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := &ValidationError{Errors: tt.errors}
			if got := ve.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidationError_Add(t *testing.T) {
	ve := NewValidationError()

	if len(ve.Errors) != 0 {
		t.Errorf("initial errors length = %v, want 0", len(ve.Errors))
	}

	ve.Add(NewFieldError("Test", "error"))

	if len(ve.Errors) != 1 {
		t.Errorf("errors length after Add = %v, want 1", len(ve.Errors))
	}
}

func TestValidationError_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   []*FieldError
		expected bool
	}{
		{
			name:     "no errors",
			errors:   []*FieldError{},
			expected: false,
		},
		{
			name: "has errors",
			errors: []*FieldError{
				{Field: "Test", Message: "error"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := &ValidationError{Errors: tt.errors}
			if got := ve.HasErrors(); got != tt.expected {
				t.Errorf("HasErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidationError_ToError(t *testing.T) {
	tests := []struct {
		name     string
		errors   []*FieldError
		wantErr  bool
	}{
		{
			name:     "no errors",
			errors:   []*FieldError{},
			wantErr:  false,
		},
		{
			name: "has errors",
			errors: []*FieldError{
				{Field: "Test", Message: "error"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := &ValidationError{Errors: tt.errors}
			err := ve.ToError()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	ve := NewValidationError()

	if ve.Errors == nil {
		t.Error("Errors slice should not be nil")
	}

	if len(ve.Errors) != 0 {
		t.Errorf("initial errors length = %v, want 0", len(ve.Errors))
	}
}

func TestProviderError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ProviderError
		expected string
	}{
		{
			name: "with message",
			err: &ProviderError{
				Message: "payment failed",
			},
			expected: "payment failed",
		},
		{
			name:     "without message",
			err:      &ProviderError{},
			expected: "provider error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProviderError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	err := &ProviderError{
		Message: "test error",
		Err:     baseErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestNewProviderError(t *testing.T) {
	raw := map[string]interface{}{
		"code":    "400",
		"details": "test",
	}

	err := NewProviderError("midtrans", "400", "Bad Request", raw)

	if err.Provider != "midtrans" {
		t.Errorf("Provider = %v, want midtrans", err.Provider)
	}

	if err.Code != "400" {
		t.Errorf("Code = %v, want 400", err.Code)
	}

	if err.Message != "Bad Request" {
		t.Errorf("Message = %v, want Bad Request", err.Message)
	}

	if err.Raw == nil {
		t.Error("Raw should not be nil")
	}
}

func TestWrapProviderError(t *testing.T) {
	raw := map[string]interface{}{
		"status": "error",
	}

	err := WrapProviderError("xendit", "500", "Internal Error", raw)

	if pe, ok := err.(*ProviderError); !ok {
		t.Error("WrapProviderError should return *ProviderError")
	} else {
		if pe.Provider != "xendit" {
			t.Errorf("Provider = %v, want xendit", pe.Provider)
		}
	}
}

func TestIsProviderError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "provider error",
			err:      &ProviderError{},
			expected: true,
		},
		{
			name:     "wrapped provider error",
			err:      &FieldError{Err: &ProviderError{}},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("test"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsProviderError(tt.err); got != tt.expected {
				t.Errorf("IsProviderError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "validation error",
			err:      &ValidationError{},
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("test"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.expected {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsFieldError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "field error",
			err:      &FieldError{},
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("test"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFieldError(tt.err); got != tt.expected {
				t.Errorf("IsFieldError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorVariables(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrUnimplemented", ErrUnimplemented, "not yet implement for this payment method"},
		{"ErrMissingParameter", ErrMissingParameter, "missing required parameter"},
		{"ErrInvalidParameter", ErrInvalidParameter, "invalid parameter"},
		{"ErrInvalidCredentials", ErrInvalidCredentials, "invalid credentials"},
		{"ErrInvalidSignature", ErrInvalidSignature, "invalid webhook signature"},
		{"ErrInvalidPayload", ErrInvalidPayload, "invalid webhook payload"},
		{"ErrMissingCredentials", ErrMissingCredentials, "missing credentials"},
		{"ErrMinAmount", ErrMinAmount, "minimum transaction amount is Rp10.000"},
		{"ErrDuplicateTransaction", ErrDuplicateTransaction, "duplicate transaction ID"},
		{"ErrTransactionNotFound", ErrTransactionNotFound, "transaction not found"},
		{"ErrTransactionFailed", ErrTransactionFailed, "transaction failed"},
		{"ErrInvalidPhoneNumber", ErrInvalidPhoneNumber, "numeric only with min length 2 or max length 13 digit. start with +62 for ID or +63 for PH"},
		{"ErrTimeout", ErrTimeout, "request timeout"},
		{"ErrRateLimit", ErrRateLimit, "rate limit exceeded"},
		{"ErrServiceUnavailable", ErrServiceUnavailable, "service unavailable"},
		{"ErrNetworkError", ErrNetworkError, "network error"},
		{"ErrWebhookVerificationFailed", ErrWebhookVerificationFailed, "webhook verification failed"},
		{"ErrInvalidWebhookType", ErrInvalidWebhookType, "invalid webhook type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
