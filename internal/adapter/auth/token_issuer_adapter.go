package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type TokenIssuerAdapter struct{}

func NewTokenIssuerAdapter() *TokenIssuerAdapter {
	return &TokenIssuerAdapter{}
}

func (ti TokenIssuerAdapter) Generate(claims port.AccessTokenPayload) (string, error) {
	mapClaims := jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   claims.ID,
		"email": claims.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("token_issuer_adapter: failed to sign token: %w", err)
	}
	return signed, nil
}

func (ti TokenIssuerAdapter) Validate(tokenStr string) (*port.AccessTokenPayload, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token_issuer_adapter: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return &port.AccessTokenPayload{}, fmt.Errorf("token_issuer_adapter: failed to parse token: %w", err)
	}

	payload, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return &port.AccessTokenPayload{}, fmt.Errorf("token_issuer_adapter: failed to cast claims to map claims")
	}

	return &port.AccessTokenPayload{
		ID:    payload["sub"].(string),
		Email: payload["email"].(string),
	}, nil
}
