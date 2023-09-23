package domain

import "errors"

var (
	ErrRequired   = errors.New("required value")
	ErrDateFormat = errors.New("invalid date request format")
)
