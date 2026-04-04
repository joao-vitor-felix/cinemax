package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHasherAdapter struct{}

func NewPasswordHasherAdapter() *PasswordHasherAdapter {
	return &PasswordHasherAdapter{}
}

func (h *PasswordHasherAdapter) Hash(password string) (string, error) {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *PasswordHasherAdapter) Compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
