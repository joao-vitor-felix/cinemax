package domain

import "time"

type RefreshToken struct {
	Token     string
	UserId    string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *RefreshToken) IsUsed() bool {
	return t.UsedAt != nil
}
