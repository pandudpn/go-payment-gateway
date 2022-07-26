package pg

import (
	"errors"
)

var (
	ErrInvalidParameter   = errors.New("parameter invalid")
	ErrNilCredentials     = errors.New("credentials nil")
	ErrMinAmount          = errors.New("minimum transaction with e-wallet is Rp100 or ₱1")
	ErrInvalidPhoneNumber = errors.New("numeric only with min length 2 or max length 13 digit. start with +62 for ID or +63 for PH")
)
