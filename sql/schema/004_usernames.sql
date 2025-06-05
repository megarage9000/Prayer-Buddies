-- +goose Up
ALTER TABLE users
ADD COLUMN username TEXT NOT NULL default 'unset';

-- +goose Down
ALTER TABLE users
DROP COLUMN username;