package port

import "time"

type TokenIssuer[C any] interface {
	Generate(claims C, expiresIn time.Duration) (string, error)
	Validate(token string) (C, error)
}

type AccessTokenPayload struct {
	Id    string
	Email string
	Role  string
}
