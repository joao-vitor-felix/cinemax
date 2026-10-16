package domain

import "time"

type PasswordResetToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

func (t *PasswordResetToken) Validate() error {
	if t.IsUsed() {
		return TokenAlreadyUsedError
	}
	if t.IsExpired() {
		return TokenExpiredError
	}
	return nil
}
