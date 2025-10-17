-- +goose Up
CREATE TYPE auth_provider AS ENUM ('credentials', 'oauth2');

ALTER TABLE users
ADD COLUMN auth_provider auth_provider NOT NULL DEFAULT 'credentials';

ALTER TABLE users ADD CONSTRAINT check_auth_provider_password_consistency CHECK (
	auth_provider = 'credentials'
	AND hashed_password IS NOT NULL
	OR auth_provider = 'oauth2'
	AND hashed_password IS NULL
);

-- +goose Down
DROP TYPE IF EXISTS auth_provider CASCADE;
