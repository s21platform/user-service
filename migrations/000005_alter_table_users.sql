-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS invite_link VARCHAR(255) NOT NULL DEFAULT 'https://space-21.ru';

-- +goose Down
ALTER TABLE users
DROP COLUMN invite_link;
