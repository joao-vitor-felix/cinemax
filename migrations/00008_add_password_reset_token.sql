-- +goose Up
CREATE TABLE IF NOT EXISTS password_reset_tokens (
	token VARCHAR(255) NOT NULL,
	user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
	expires_at TIMESTAMPTZ NOT NULL,
	used_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (token, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS password_reset_tokens;
