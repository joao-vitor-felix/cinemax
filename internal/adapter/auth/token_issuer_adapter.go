package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenIssuerAdapter[C any] struct {
	secret string
}

func NewTokenIssuerAdapter[C any](secret string) *TokenIssuerAdapter[C] {
	return &TokenIssuerAdapter[C]{secret: secret}
}

func (ti *TokenIssuerAdapter[C]) Generate(claims C, expiresIn time.Duration) (string, error) {
	data, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("token_issuer_adapter: failed to marshal claims: %w", err)
	}

	mapClaims := jwt.MapClaims{
		"exp": time.Now().Add(expiresIn).Unix(),
		"iat": time.Now().Unix(),
	}
	if err := json.Unmarshal(data, &mapClaims); err != nil {
		return "", fmt.Errorf("token_issuer_adapter: failed to build map claims: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	signed, err := token.SignedString([]byte(ti.secret))
	if err != nil {
		return "", fmt.Errorf("token_issuer_adapter: failed to sign token: %w", err)
	}

	return signed, nil
}

func (ti *TokenIssuerAdapter[C]) Validate(tokenStr string) (C, error) {
	payload := *new(C)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token_issuer_adapter: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ti.secret), nil
	})

	if err != nil || !token.Valid {
		return payload, fmt.Errorf("token_issuer_adapter: failed to parse token: %w", err)
	}

	data, err := json.Marshal(token.Claims)
	if err != nil {
		return payload, fmt.Errorf("token_issuer_adapter: failed to marshal map claims: %w", err)
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, fmt.Errorf("token_issuer_adapter: failed to unmarshal claims: %w", err)
	}

	return payload, nil
}
