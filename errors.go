package pg

import (
	"errors"
	"fmt"
)

var (
	// Base errors
	ErrUnimplemented       = errors.New("not yet implement for this payment method")
	ErrMissingParameter    = errors.New("missing required parameter")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidSignature    = errors.New("invalid webhook signature")
	ErrInvalidPayload      = errors.New("invalid webhook payload")
	ErrMissingCredentials  = errors.New("missing credentials")

	// Payment errors
	ErrMinAmount          = errors.New("minimum transaction amount is Rp10.000")
	ErrDuplicateTransaction = errors.New("duplicate transaction ID")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrTransactionFailed  = errors.New("transaction failed")
	ErrInvalidPhoneNumber = errors.New("numeric only with min length 2 or max length 13 digit. start with +62 for ID or +63 for PH")

	// Network errors
	ErrTimeout            = errors.New("request timeout")
	ErrRateLimit          = errors.New("rate limit exceeded")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrNetworkError       = errors.New("network error")

	// Webhook errors
	ErrWebhookVerificationFailed = errors.New("webhook verification failed")
	ErrInvalidWebhookType        = errors.New("invalid webhook type")
)

// FieldError represents an error for a specific field
type FieldError struct {
	Field   string
	Message string
	Err     error
}

// Error returns the error message
func (e *FieldError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("%s is invalid", e.Field)
}

// Unwrap returns the underlying error
func (e *FieldError) Unwrap() error {
	return e.Err
}

// NewFieldError creates a new FieldError
func NewFieldError(field, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
		Err:     ErrInvalidParameter,
	}
}

// NewRequiredFieldError creates a new FieldError for missing required field
func NewRequiredFieldError(field string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: "is required",
		Err:     ErrMissingParameter,
	}
}

// ValidationError represents a collection of field errors
type ValidationError struct {
	Errors []*FieldError
}

// Error returns the error message
func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("validation failed: %d errors", len(e.Errors))
}

// Add adds a field error to the validation error
func (e *ValidationError) Add(err *FieldError) {
	e.Errors = append(e.Errors, err)
}

// HasErrors returns true if there are any errors
func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ToError returns the validation error or nil if there are no errors
func (e *ValidationError) ToError() error {
	if e.HasErrors() {
		return e
	}
	return nil
}

// NewValidationError creates a new ValidationError
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make([]*FieldError, 0),
	}
}

// WrapProviderError wraps a provider-specific error into a unified error
func WrapProviderError(provider, code, message string, raw map[string]interface{}) error {
	return &ProviderError{
		Code:     code,
		Message:  message,
		Provider: provider,
		Raw:      raw,
	}
}

// IsProviderError checks if an error is a ProviderError
func IsProviderError(err error) bool {
	_, ok := err.(*ProviderError)
	return ok
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsFieldError checks if an error is a FieldError
func IsFieldError(err error) bool {
	_, ok := err.(*FieldError)
	return ok
}
