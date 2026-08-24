-- +goose Up
CREATE TABLE IF NOT EXISTS game_users (
    user_id UUID PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS game_users;
