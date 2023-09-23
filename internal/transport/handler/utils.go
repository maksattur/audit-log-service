package handler

import (
	"github.com/maksattur/audit-log-service/internal"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	if password == "" {
		return "", ErrPasswordIsEmpty
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &internal.CustomError{
			OriginalError: err,
			Message:       "generate password",
		}
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
