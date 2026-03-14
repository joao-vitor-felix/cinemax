package port

type TokenIssuer interface {
	Generate(claims AccessTokenPayload) (string, error)
	Validate(token string) (*AccessTokenPayload, error)
}

type AccessTokenPayload struct {
	Id    string
	Email string
	Role  string
}
