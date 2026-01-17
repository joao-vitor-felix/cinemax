-- +goose Up
ALTER TABLE users
DROP CONSTRAINT check_auth_provider_password_consistency;

DROP TYPE IF EXISTS auth_provider CASCADE;

-- +goose Down
CREATE TYPE auth_provider AS ENUM ('credentials', 'oauth2');

ALTER TABLE users
ADD COLUMN auth_provider auth_provider NOT NULL DEFAULT 'credentials',
ADD CONSTRAINT check_auth_provider_password_consistency CHECK (
	auth_provider = 'credentials'
	AND password_hash IS NOT NULL
	OR auth_provider = 'oauth2'
	AND password_hash IS NULL
);
