-- +goose Up

CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    initiator UUID,
    invited UUID,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS friends;