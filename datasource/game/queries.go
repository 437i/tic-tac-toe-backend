package game

const LeaderboardQuery = `
WITH player_games AS (
    SELECT
        player_x AS user_id,
        CASE WHEN status = 'PlayerXWon' THEN 1 ELSE 0 END AS win,
        CASE WHEN status = 'PlayerOWon' THEN 1 ELSE 0 END AS loss,
        CASE WHEN status = 'Draw' THEN 1 ELSE 0 END AS draw
    FROM
        games
    WHERE status IN (
        'PlayerXWon',
        'PlayerOWon',
        'Draw'
    )

    UNION ALL

    SELECT
        player_o AS user_id,
        CASE WHEN status = 'PlayerOWon' THEN 1 ELSE 0 END AS win,
        CASE WHEN status = 'PlayerXWon' THEN 1 ELSE 0 END AS loss,
        CASE WHEN status = 'Draw' THEN 1 ELSE 0 END AS draw
    FROM games
    WHERE status IN (
        'PlayerXWon',
        'PlayerOWon',
        'Draw'
    )
),
stats AS (
    SELECT
        user_id,
        COUNT(*) AS total_games,
        SUM(win) AS wins,
        SUM(loss) AS losses,
        SUM(draw) AS draws
    FROM player_games
    WHERE user_id IS NOT NULL
    GROUP BY user_id
)

SELECT
    s.user_id,
    gu.login,
    s.total_games,
    s.wins,
    s.losses,
    s.draws,
    CASE
        WHEN s.losses + s.draws = 0 THEN s.wins::NUMERIC
        ELSE ROUND(
                s.wins::NUMERIC / (s.losses + s.draws),
                4
            )
    END AS winrate
FROM stats s
JOIN game_users gu
    ON s.user_id = gu.user_id
ORDER BY
    winrate DESC NULLS LAST,
    s.wins DESC,
    s.total_games DESC,
    gu.login ASC
LIMIT $1;
`
