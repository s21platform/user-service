-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS posts (
    id          serial PRIMARY KEY,
    user_id     integer NOT NULL,
    content     TEXT,
    external_id TEXT,
    origin      TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS posts;
DROP EXTENSION IF EXISTS pgcrypto;