-- +goose Up
CREATE TYPE user_gender AS ENUM ('male', 'female', 'other', 'prefer_not_to_say');

CREATE TABLE
	IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(254) UNIQUE NOT NULL,
		phone VARCHAR(15) UNIQUE,
		hashed_password VARCHAR(255),
		-- it can be used to verify age-restricted content and suggest content based on age
		date_of_birth DATE NOT NULL,
		gender user_gender,
		profile_photo_url VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
  BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE OR REPLACE TRIGGER set_updated_at_trigger
AFTER UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TYPE IF EXISTS user_gender;
DROP FUNCTION IF EXISTS set_updated_at() CASCADE;
DROP TABLE IF EXISTS users;
