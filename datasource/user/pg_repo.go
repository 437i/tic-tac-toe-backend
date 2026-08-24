package user

import (
	domainUser "apg105/domain/user"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PGRepo) Create(ctx context.Context, userDm domainUser.User) error {
	user := toRepoUser(userDm)
	sql, args, err := r.queryBuilder.
		Insert(usersTable).
		Columns(colUserID, colLogin, colPassword).
		Values(user.UserID, user.Login, user.Password).
		ToSql()
	if err != nil {
		return fmt.Errorf("building create query: %w", err)
	}
	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return domainUser.ErrUserAlreadyExists
		}
		return fmt.Errorf("executing create query: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", result.RowsAffected())
	}
	return nil
}

func (r *PGRepo) GetUserByLogin(ctx context.Context, login string) (domainUser.User, error) {
	sql, args, err := r.queryBuilder.
		Select(colUserID, colLogin, colPassword).
		From(usersTable).
		Where(sq.Eq{colLogin: login}).
		ToSql()
	if err != nil {
		return domainUser.User{}, fmt.Errorf("building get query: %w", err)
	}
	var user User
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.UserID,
		&user.Login,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainUser.User{}, domainUser.ErrUserNotFound
		}
		return domainUser.User{}, fmt.Errorf("executing get query: %w", err)
	}
	return user.toDomain(), nil
}

func (r *PGRepo) GetLoginByID(ctx context.Context, userID uuid.UUID) (string, error) {
	sql, args, err := r.queryBuilder.
		Select(colLogin).
		From(usersTable).
		Where(sq.Eq{colUserID: userID}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("building get query: %w", err)
	}
	var login string
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domainUser.ErrUserNotFound
		}
		return "", fmt.Errorf("executing get query: %w", err)
	}
	return login, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
