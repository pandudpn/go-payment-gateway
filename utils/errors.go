package utils

import (
	"errors"
)

var (
	ErrHttpRequest       = errors.New("http request is nil")
	ErrHttpClient        = errors.New("http client is nil")
	ErrParseBodyResponse = errors.New("error parse body response")
)
