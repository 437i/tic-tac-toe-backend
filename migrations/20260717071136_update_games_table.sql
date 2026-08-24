-- +goose Up
ALTER TABLE games
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE games
    ADD CONSTRAINT fk_player_x
        FOREIGN KEY (player_x)
        REFERENCES game_users(user_id)
        ON DELETE RESTRICT;

ALTER TABLE games
    ADD CONSTRAINT fk_player_o
        FOREIGN KEY (player_o)
        REFERENCES game_users(user_id)
        ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_player_x
    ON games(player_x);

CREATE INDEX IF NOT EXISTS idx_player_o
    ON games(player_o);

CREATE INDEX IF NOT EXISTS idx_status_created_at
    ON games(status, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_status_created_at;
DROP INDEX IF EXISTS idx_player_x;
DROP INDEX IF EXISTS idx_player_o;

ALTER TABLE games
    DROP CONSTRAINT IF EXISTS fk_player_x;

ALTER TABLE games
    DROP CONSTRAINT IF EXISTS fk_player_o;

ALTER TABLE games
    DROP COLUMN IF EXISTS created_at;