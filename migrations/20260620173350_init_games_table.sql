-- +goose Up
CREATE TABLE IF NOT EXISTS games (
    game_id UUID PRIMARY KEY,

    field INT[9] NOT NULL
        CHECK (
            cardinality(field) = 9
            AND array_position(field, NULL) IS NULL
            AND field <@ ARRAY[0, 1, 2]::INT[]
        ),

    player_x UUID,
    player_o UUID,

    status TEXT NOT NULL
        CHECK (
            status IN (
                'WaitingForOpponent',
                'PlayerXTurn',
                'PlayerOTurn',
                'PlayerXWon',
                'PlayerOWon',
                'Draw'
            )
        ),

    mode TEXT NOT NULL
        CHECK (
            mode IN ('PvP', 'PvE')
        ),

    version INTEGER NOT NULL DEFAULT 1
        CHECK (version > 0),

    CONSTRAINT chk_games_players_distinct
        CHECK (
            player_x IS DISTINCT FROM player_o
        ),

    CONSTRAINT chk_games_mode_state
        CHECK (
            (
                mode = 'PvE'
                AND status != 'WaitingForOpponent'
                AND (
                    (player_x IS NOT NULL AND player_o IS NULL)
                    OR
                    (player_x IS NULL AND player_o IS NOT NULL)
                )
            )
            OR
            (
                mode = 'PvP'
                AND (
                    (
                        status = 'WaitingForOpponent'
                        AND (
                            (player_x IS NOT NULL AND player_o IS NULL)
                            OR
                            (player_x IS NULL AND player_o IS NOT NULL)
                        )
                    )
                    OR
                    (
                        status IN (
                            'PlayerXTurn',
                            'PlayerOTurn',
                            'PlayerXWon',
                            'PlayerOWon',
                            'Draw'
                        )
                        AND player_x IS NOT NULL
                        AND player_o IS NOT NULL
                    )
                )
            )
        )
);

-- +goose Down
DROP TABLE IF EXISTS games;
