package port

type PasswordHasher interface {
	Hash(password []byte) ([]byte, error)
	Compare(password, hash []byte) error
}
