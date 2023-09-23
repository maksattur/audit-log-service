package token_manager

import "errors"

var (
	ErrTokenIsExpired = errors.New("token is expired")
)
