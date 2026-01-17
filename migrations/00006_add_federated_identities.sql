-- +goose Up
CREATE TABLE IF NOT EXISTS federated_identities (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		provider VARCHAR(20) NOT NULL,
		provider_user_id VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE(user_id, provider),
		UNIQUE(provider, provider_user_id)
);

-- +goose Down
DROP TABLE IF EXISTS federated_identities;
