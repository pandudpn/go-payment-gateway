package pg

import (
	"errors"
)

var (
	ErrInvalidParameter = errors.New("parameter invalid")
	ErrNilCredentials   = errors.New("credentials nil")
)
