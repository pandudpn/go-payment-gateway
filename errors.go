package pg

import (
	"errors"
)

var (
	ErrUnimplemented      = errors.New("not yet implement for this payment method")
	ErrMissingParameter   = errors.New("missing parameter")
	ErrInvalidParameter   = errors.New("parameter invalid")
	ErrMissingCredentials = errors.New("missing credentials")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrMinAmount          = errors.New("minimum transaction with e-wallet is Rp100 or ₱1")
	ErrInvalidPhoneNumber = errors.New("numeric only with min length 2 or max length 13 digit. start with +62 for ID or +63 for PH")
)
