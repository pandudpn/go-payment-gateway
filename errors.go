package pg

import (
	"errors"
)

var (
	ErrInvalidParameter = errors.New("parameter invalid")
	ErrNilCredentials   = errors.New("credentials nil")
	ErrMinAmount        = errors.New("minimum transaction with e-wallet is Rp100 or ₱1")
)
