package port

type TokenIssuer interface {
	Generate(claims AccessTokenPayload) (string, error)
	Validate(token string) (*AccessTokenPayload, error)
}

type AccessTokenPayload struct {
	ID    string
	Email string
}
