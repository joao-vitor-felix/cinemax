-- +goose Up
ALTER TABLE users
ALTER COLUMN id SET DEFAULT uuidv7();

ALTER TABLE federated_identities
ALTER COLUMN id SET DEFAULT uuidv7();

ALTER TABLE refresh_tokens
ALTER COLUMN token SET DEFAULT uuidv7();

-- +goose Down
ALTER TABLE refresh_tokens
ALTER COLUMN token DROP DEFAULT;

ALTER TABLE federated_identities
ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE users
ALTER COLUMN id SET DEFAULT gen_random_uuid();
