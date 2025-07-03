-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS posts (
    id          serial PRIMARY KEY,
    user_uuid   UUID NOT NULL,
    content     TEXT NOT NULL,
    external_id TEXT,
    origin      TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP,
    CONSTRAINT fk_post_user_uuid FOREIGN KEY (user_uuid) REFERENCES users(uuid);
);

-- +goose Down
DROP TABLE IF EXISTS posts;
DROP EXTENSION IF EXISTS pgcrypto;