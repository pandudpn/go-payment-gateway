package xendit

import (
	"errors"
	"reflect"
	
	pg "github.com/pandudpn/go-payment-gateway"
)

// ErrorCode common error in e-wallet payments
type ErrorCode error

var (
	// ErrUnknown error not identified
	ErrUnknown ErrorCode = errors.New("error unknown. expected error")

	// ErrAPIValidation there's invalid input in one of the required request fields
	ErrAPIValidation ErrorCode = errors.New("one or more parameters is invalid")

	// ErrEWalletNotSupported your requested e-wallet_type is not supported yet
	ErrEWalletNotSupported ErrorCode = errors.New("e-wallet not yet support")

	// ErrDuplicateTransaction the payment with the same ref_id has already made before
	ErrDuplicateTransaction ErrorCode = errors.New("duplicate reference_id")

	// ErrRequestForbidden api key in use does not have necessary permissions
	ErrRequestForbidden ErrorCode = errors.New("you don't have permission to access this API")

	// ErrUserNotAuthorized user didn't authorize the payment request within the time limit
	ErrUserNotAuthorized ErrorCode = errors.New("user don't have authorized to request the payment")

	// ErrAccountBlocked user account is blocked
	ErrAccountBlocked ErrorCode = errors.New("user account is blocked")

	// ErrSendTransaction while sending transaction notification to OVO
	ErrSendTransaction ErrorCode = errors.New("error send transaction notification to OVO")

	// ErrUserDeclined the payment request
	ErrUserDeclined ErrorCode = errors.New("user declined the payment request")

	// ErrPhoneNumberNotRegistered phone number is not yet registered in the e-wallet
	ErrPhoneNumberNotRegistered ErrorCode = errors.New("phone number not yet registered in the e-wallet system")

	// ErrEWalletUnreachable e-wallet system cannot reach the user e-wallet app/phone.
	ErrEWalletUnreachable ErrorCode = errors.New("e-wallet system can't reach customer app")

	// ErrOvoTimeout connection timeout from OVO app to the OVO server/system
	ErrOvoTimeout ErrorCode = errors.New("connection timeout to ovo server")

	// ErrCredentials merchant is not registered in e-wallet system
	ErrCredentials ErrorCode = errors.New("merchant not yet registered in e-wallet system")

	// ErrAccountAuthentication failed to authenticated
	ErrAccountAuthentication ErrorCode = errors.New("authentication has failed")

	// ErrExternalError error in e-wallet system
	ErrExternalError ErrorCode = errors.New("there's an error on e-wallet system")

	// ErrUnsupportedCurrency currency of payment not yet supported
	ErrUnsupportedCurrency ErrorCode = errors.New("the payment currency request is not support for this channel")

	// ErrInvalidPaymentMethodID mismatch between the requested channel_code or customer_id
	ErrInvalidPaymentMethodID ErrorCode = errors.New("channel_code or customer_id not match")

	// ErrInvalidAPIKey api key is invalid
	ErrInvalidAPIKey ErrorCode = errors.New("api key is invalid")

	// ErrInvalidToken error when account linking is expired
	ErrInvalidToken ErrorCode = errors.New("account linked is expired")

	// ErrChannelNotActivated channel not yet active
	ErrChannelNotActivated ErrorCode = errors.New("channel_code is not yet active")

	// ErrCallbackNotFound callback not found in request
	ErrCallbackNotFound ErrorCode = errors.New("callback url not found in request")

	// ErrContentType error when content type request not supported
	ErrContentType ErrorCode = errors.New("content type is not yet supported")

	// ErrServerError error when server gateway has unexpected error occurred
	ErrServerError ErrorCode = errors.New("gateway has unexpected error")

	// ErrChannelUnavailable error when channel is not available
	ErrChannelUnavailable ErrorCode = errors.New("payment channel is unavailable")
)

// errorCodeMap mapping an error from response into interface error
var errorCodeMap = map[string]ErrorCode{
	"API_VALIDATION_ERROR":            ErrAPIValidation,
	"EWALLET_TYPE_NOT_SUPPORTED":      ErrEWalletNotSupported,
	"DUPLICATE_PAYMENT_REQUEST_ERROR": ErrDuplicateTransaction,
	"DUPLICATE_PAYMENT":               ErrDuplicateTransaction,
	"REQUEST_FORBIDDEN_ERROR":         ErrRequestForbidden,
	"USER_DID_NOT_AUTHORIZED":         ErrUserNotAuthorized,
	"ACCOUNT_BLOCKED_ERROR":           ErrAccountBlocked,
	"SENDING_TRANSACTION_ERROR":       ErrSendTransaction,
	"USER_DECLINED":                   ErrUserDeclined,
	"PHONE_NUMBER_NOT_REGISTERED":     ErrPhoneNumberNotRegistered,
	"EWALLET_APP_UNREACHABLE":         ErrEWalletUnreachable,
	"OVO_TIMEOUT_ERROR":               ErrOvoTimeout,
	"CREDENTIALS_ERROR":               ErrCredentials,
	"ACCOUNT_AUTHENTICATION_ERROR":    ErrAccountAuthentication,
	"EXTERNAL_ERROR":                  ErrExternalError,
	"UNSUPPORTED_CURRENCY":            ErrUnsupportedCurrency,
	"INVALID_PAYMENT_METHOD_ID":       ErrInvalidPaymentMethodID,
	"INVALID_API_KEY":                 ErrInvalidAPIKey,
	"INVALID_MERCHANT_CREDENTIALS":    ErrCredentials,
	"INVALID_TOKEN":                   ErrInvalidToken,
	"CHANNEL_NOT_ACTIVATED":           ErrChannelNotActivated,
	"CALLBACK_URL_NOT_FOUND":          ErrCallbackNotFound,
	"UNSUPPORTED_CONTENT_TYPE":        ErrContentType,
	"SERVER_ERROR":                    ErrServerError,
	"CHANNEL_UNAVAILABLE":             ErrChannelUnavailable,
}

// GetErrorCode getting error code from response into interface error
func GetErrorCode(param interface{}) ErrorCode {
	var (
		er  string
		msg string
	)
	switch p := param.(type) {
	case VirtualAccount:
		er = p.ErrorCode
		msg = p.Message
	case ChargeResponse:
		er = p.ErrorCode
		msg = p.Message
	}

	if reflect.ValueOf(er).IsZero() {
		return nil
	}

	pg.Log.Error(msg)
	if v, ok := errorCodeMap[er]; ok {
		return v
	}

	return ErrUnknown
}
