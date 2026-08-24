package game

import (
	domainGame "apg105/domain/game"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepo struct {
	queryBuilder sq.StatementBuilderType
	pool         *pgxpool.Pool
}

func NewPGRepo(pool *pgxpool.Pool) *PGRepo {
	return &PGRepo{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
	}
}

func (r *PGRepo) Create(ctx context.Context, gameDm domainGame.Game) error {
	game := toRepoGame(gameDm)

	sql, args, err := r.queryBuilder.
		Insert(gamesTable).
		Columns(
			colGameID,
			colField,
			colPlayerX,
			colPlayerO,
			colStatus,
			colMode,
		).
		Values(
			game.GameId,
			game.Field,
			game.PlayerX,
			game.PlayerO,
			game.Status,
			game.Mode,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create game query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute create game query: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", result.RowsAffected())
	}

	return nil
}

func (r *PGRepo) Get(ctx context.Context, gameId uuid.UUID) (domainGame.Game, error) {
	sql, args, err := r.queryBuilder.
		Select(
			colGameID,
			colField,
			colPlayerX,
			colPlayerO,
			colStatus,
			colMode,
			colVersion,
			colCreatedAt,
		).
		From(gamesTable).
		Where(sq.Eq{colGameID: gameId}).
		ToSql()
	if err != nil {
		return domainGame.Game{}, fmt.Errorf("build get game query: %w", err)
	}

	var row Game

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&row.GameId,
		&row.Field,
		&row.PlayerX,
		&row.PlayerO,
		&row.Status,
		&row.Mode,
		&row.Version,
		&row.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainGame.Game{}, domainGame.ErrGameNotFound
		}
		return domainGame.Game{}, fmt.Errorf("execute get game query: %w", err)
	}

	game, err := row.toDomain()
	if err != nil {
		return domainGame.Game{}, fmt.Errorf("map game to domain: %w", err)
	}

	return game, nil
}

func (r *PGRepo) GetAvailable(ctx context.Context, userId uuid.UUID) ([]domainGame.GameMeta, error) {
	sql, args, err := r.queryBuilder.
		Select(
			colGameID,
			colPlayerX,
			colPlayerO,
			colStatus,
			colMode,
			colCreatedAt,
		).
		From(gamesTable).
		Where(sq.Eq{colStatus: statusWaitingForOpponent}).
		Where(sq.Expr(
			colPlayerX+" IS DISTINCT FROM ?",
			userId,
		)).
		Where(sq.Expr(
			colPlayerO+" IS DISTINCT FROM ?",
			userId,
		)).
		OrderBy(colCreatedAt + " ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get available query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute get available query: %w", err)
	}
	defer rows.Close()

	games, err := scanGamesMeta(rows)
	if err != nil {
		return nil, fmt.Errorf("scan game meta: %w", err)
	}

	return games, nil
}

func (r *PGRepo) GetFinished(ctx context.Context, userId uuid.UUID) ([]domainGame.GameMeta, error) {
	sql, args, err := r.queryBuilder.
		Select(
			colGameID,
			colPlayerX,
			colPlayerO,
			colStatus,
			colMode,
			colCreatedAt,
		).
		From(gamesTable).
		Where(sq.Or{
			sq.Eq{colPlayerX: userId},
			sq.Eq{colPlayerO: userId},
		}).
		Where(sq.Eq{colStatus: []string{
			statusPlayerXWon,
			statusPlayerOWon,
			statusDraw,
		}}).
		OrderBy(colCreatedAt + " DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get finished query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute get finished query: %w", err)
	}

	defer rows.Close()

	games, err := scanGamesMeta(rows)
	if err != nil {
		return nil, fmt.Errorf("scan game meta: %w", err)
	}

	return games, nil
}

func (r *PGRepo) Update(ctx context.Context, gameDm domainGame.Game) (int, error) {
	game := toRepoGame(gameDm)

	sql, args, err := r.queryBuilder.
		Update(gamesTable).
		Set(colField, game.Field).
		Set(colStatus, game.Status).
		Set(colPlayerX, game.PlayerX).
		Set(colPlayerO, game.PlayerO).
		Set(colVersion, sq.Expr("version + 1")).
		Where(sq.Eq{colGameID: game.GameId}).
		Where(sq.Eq{colVersion: game.Version}).
		Suffix("RETURNING " + colVersion).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build update query: %w", err)
	}

	var newVersion int

	err = r.pool.QueryRow(ctx, sql, args...).Scan(&newVersion)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			exists, existsErr := r.gameExists(ctx, game.GameId)
			if existsErr != nil {
				return 0, fmt.Errorf("check game exists: %w", existsErr)
			}
			if !exists {
				return 0, domainGame.ErrGameNotFound
			}
			return 0, domainGame.ErrConcurrentModification
		}
		return 0, fmt.Errorf("execute update query: %w", err)
	}

	return newVersion, nil
}

func (r *PGRepo) GetLeaderboard(ctx context.Context, limit int) ([]domainGame.UserStats, error) {
	rows, err := r.pool.Query(ctx, LeaderboardQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("executing get leaderboard query: %w", err)
	}

	defer rows.Close()

	leaders := make([]domainGame.UserStats, 0, limit)

	for rows.Next() {
		row := UserStats{}
		if err := rows.Scan(
			&row.UserID,
			&row.Login,
			&row.TotalGames,
			&row.Wins,
			&row.Losses,
			&row.Draws,
			&row.Winrate,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		leaders = append(leaders, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return leaders, nil
}

func scanGamesMeta(rows pgx.Rows) ([]domainGame.GameMeta, error) {
	games := make([]domainGame.GameMeta, 0)

	for rows.Next() {
		row := GameMeta{}

		if err := rows.Scan(
			&row.GameID,
			&row.PlayerX,
			&row.PlayerO,
			&row.Status,
			&row.Mode,
			&row.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		game, err := row.toDomain()
		if err != nil {
			return nil, fmt.Errorf("map game meta to domain: %w", err)
		}

		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return games, nil
}

func (r *PGRepo) gameExists(ctx context.Context, gameID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS (
		SELECT 1
		FROM games
		WHERE game_id = $1
	)`

	var exists bool

	if err := r.pool.QueryRow(ctx, query, gameID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
