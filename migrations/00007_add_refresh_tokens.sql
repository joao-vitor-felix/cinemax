-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
	token VARCHAR(255) PRIMARY KEY,
	user_id UUID REFERENCES users,
	expires_at TIMESTAMPTZ NOT NULL,
	used_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN refresh_tokens.token IS 'A SHA-256 hash of token';

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
