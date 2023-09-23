package handler

import (
	"errors"
)

var (
	ErrLoginOrPasswordIsEmpty   = errors.New("login or password is empty")
	ErrLoginOrPasswordIncorrect = errors.New("login or password is incorrect")

	ErrPasswordIsEmpty = errors.New("password is empty")
)
