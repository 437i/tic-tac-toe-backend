package game

import (
	"context"

	"github.com/google/uuid"
)

type AI interface {
	Move(ctx context.Context, game Game) (Game, error)
}

type Repo interface {
	Create(ctx context.Context, gameDm Game) error
	Get(ctx context.Context, gameId uuid.UUID) (Game, error)
	GetAvailable(ctx context.Context, userId uuid.UUID) ([]GameMeta, error)
	GetFinished(ctx context.Context, userId uuid.UUID) ([]GameMeta, error)
	Update(ctx context.Context, gameDm Game) (int, error)
	GetLeaderboard(ctx context.Context, limit int) ([]UserStats, error)
}
