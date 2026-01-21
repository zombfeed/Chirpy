-- +goose up
ALTER TABLE users
ADD password TEXT NOT NULL;

-- +goose down
ALTER TABLE users
DROP COLUMN password;
