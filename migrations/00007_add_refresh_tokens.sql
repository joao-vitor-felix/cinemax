-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
	token UUID PRIMARY KEY,
	user_id UUID REFERENCES users,
	expires_at TIMESTAMPTZ NOT NULL,
	used_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
