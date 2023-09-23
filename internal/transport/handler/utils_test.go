package handler

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "qwerty123"
	hash, err := GeneratePasswordHash(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestGeneratePasswordHashError(t *testing.T) {
	password := ""
	hash, err := GeneratePasswordHash(password)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrPasswordIsEmpty)
	require.Empty(t, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "qwerty"
	hash, _ := GeneratePasswordHash(password)
	result := CheckPasswordHash(password, hash)
	require.True(t, result)
}

func TestCheckPasswordHashError(t *testing.T) {
	password := "qwerty"
	hash := "$2a$10$wRI1uRkdb3uxw9dJyHrudeSGFPPo5aIFO4LanU.GAq0YfknryquFW"
	result := CheckPasswordHash(password, hash)
	require.False(t, result)
}
