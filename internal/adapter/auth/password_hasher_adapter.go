package auth

import (
	"github.com/alexedwards/argon2id"
)

type PasswordHasherAdapter struct{}

func NewPasswordHasherAdapter() *PasswordHasherAdapter {
	return &PasswordHasherAdapter{}
}

func (h *PasswordHasherAdapter) Hash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (h *PasswordHasherAdapter) Compare(hash *string, password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, *hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
