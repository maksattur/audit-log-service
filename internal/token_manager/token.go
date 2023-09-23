package token_manager

import (
	"encoding/json"
	"github.com/cristalhq/jwt/v5"
	"time"
)

type TokenManager struct {
	builder *jwt.Builder
	key     []byte
	jwtTTL  time.Duration
}

func NewTokenManager(key []byte, jwtTTL time.Duration) (*TokenManager, error) {
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, err
	}
	b := jwt.NewBuilder(signer)
	return &TokenManager{
		builder: b,
		key:     key,
		jwtTTL:  jwtTTL,
	}, nil
}

func (tm *TokenManager) BuildToken() (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.jwtTTL)),
	}
	token, err := tm.builder.Build(claims)
	if err != nil {
		return "", err
	}
	return string(token.Bytes()), nil
}

func (tm *TokenManager) VerifyToken(tokenString string) error {
	verifier, err := jwt.NewVerifierHS(jwt.HS256, tm.key)
	if err != nil {
		return err
	}

	newToken, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		return err
	}

	var newClaims jwt.RegisteredClaims

	if err := json.Unmarshal(newToken.Claims(), &newClaims); err != nil {
		return err
	}

	if time.Unix(newClaims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return ErrTokenIsExpired
	}

	return nil
}
