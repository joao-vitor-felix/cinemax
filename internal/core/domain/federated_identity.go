package domain

import (
	"time"

	"github.com/google/uuid"
)

type FederatedIdentity struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	CreatedAt      time.Time
}
