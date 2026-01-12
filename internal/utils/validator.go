package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/pandudpn/go-payment-gateway"
)

// emailRegex is the regex pattern for email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// phoneRegex is the regex pattern for phone number validation (Indonesia format)
var phoneRegex = regexp.MustCompile(`^(\+62|62)[0-9]{9,13}$`)

// RequiredField validates that a field is not zero/empty
func RequiredField(value interface{}, fieldName string) error {
	if reflect.ValueOf(value).IsZero() {
		return pg.NewRequiredFieldError(fieldName)
	}
	return nil
}

// RequiredString validates that a string is not empty
func RequiredString(value, fieldName string) error {
	if value == "" {
		return pg.NewRequiredFieldError(fieldName)
	}
	return nil
}

// MinAmount validates minimum amount based on payment type
func MinAmount(amount int64, paymentType pg.PaymentType) error {
	var min int64

	switch {
	case paymentType.IsEWallet():
		min = pg.MinAmountEWallet
	case paymentType.IsVirtualAccount():
		min = pg.MinAmountVA
	case paymentType == pg.PaymentTypeQRIS:
		min = pg.MinAmountQRIS
	case paymentType == pg.PaymentTypeCC:
		min = pg.MinAmountCC
	case paymentType.IsRetail():
		min = pg.MinAmountRetail
	default:
		min = pg.MinAmountEWallet
	}

	if amount < min {
		return &pg.FieldError{
			Field:   "Amount",
			Message: fmt.Sprintf("must be at least Rp%d", min),
			Err:     pg.ErrMinAmount,
		}
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email, fieldName string) error {
	if email == "" {
		return pg.NewRequiredFieldError(fieldName)
	}
	if !emailRegex.MatchString(email) {
		return &pg.FieldError{
			Field:   fieldName,
			Message: "must be a valid email address",
			Err:     pg.ErrInvalidParameter,
		}
	}
	return nil
}

// ValidatePhone validates phone number for Indonesia format
// Accepts formats: +62812345678, 62812345678, 0812345678
func ValidatePhone(phone, fieldName string) error {
	if phone == "" {
		return pg.NewRequiredFieldError(fieldName)
	}

	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(strings.ReplaceAll(phone, " ", ""), "-", "")

	// Convert 08... to +628...
	if strings.HasPrefix(cleaned, "08") {
		cleaned = "+62" + cleaned[1:]
	}

	// Check if matches regex
	if !phoneRegex.MatchString(cleaned) {
		return &pg.FieldError{
			Field:   fieldName,
			Message: "must be a valid Indonesian phone number (e.g., +628123456789)",
			Err:     pg.ErrInvalidPhoneNumber,
		}
	}

	return nil
}

// InArray checks if an item exists in a list using generics
func InArray[T comparable](item T, list []T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// ValidateEnum validates that a value is in the allowed list
func ValidateEnum[T comparable](value T, allowed []T, name string) error {
	if !InArray(value, allowed) {
		var values []string
		for _, v := range allowed {
			values = append(values, fmt.Sprintf("%v", v))
		}
		return &pg.FieldError{
			Field:   name,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(values, ", ")),
			Err:     pg.ErrInvalidParameter,
		}
	}
	return nil
}

// ValidatePaymentType validates that a payment type is supported
func ValidatePaymentType(paymentType pg.PaymentType) error {
	allowed := []pg.PaymentType{
		pg.PaymentTypeGoPay, pg.PaymentTypeOVO, pg.PaymentTypeDANA,
		pg.PaymentTypeShopeePay, pg.PaymentTypeLinkAja,
		pg.PaymentTypeVABCA, pg.PaymentTypeVABNI, pg.PaymentTypeVABRI,
		pg.PaymentTypeVAMandiri, pg.PaymentTypeVAPermata, pg.PaymentTypeVACIMB,
		pg.PaymentTypeQRIS, pg.PaymentTypeCC,
		pg.PaymentTypeAlfamart, pg.PaymentTypeIndomaret,
	}
	return ValidateEnum(paymentType, allowed, "PaymentType")
}

// ValidateOrderID validates order ID format
func ValidateOrderID(orderID string) error {
	if err := RequiredString(orderID, "OrderID"); err != nil {
		return err
	}
	if len(orderID) > 100 {
		return &pg.FieldError{
			Field:   "OrderID",
			Message: "must be less than 100 characters",
			Err:     pg.ErrInvalidParameter,
		}
	}
	return nil
}

// SetDefault sets a default value if the current value is zero
func SetDefault[T any](value *T, defaultValue T) {
	if reflect.ValueOf(*value).IsZero() {
		*value = defaultValue
	}
}

// NewValidationError creates a new validation error
func NewValidationError() *pg.ValidationError {
	return pg.NewValidationError()
}
