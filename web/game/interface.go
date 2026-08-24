package game

import (
	dmg "apg105/domain/game"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateGame(ctx context.Context, req dmg.GameCreation) (dmg.Game, error)
	MakeMove(ctx context.Context, userID uuid.UUID, game dmg.Game) (dmg.Game, error)
	Join(ctx context.Context, gameId, userId uuid.UUID) error
	GetGame(ctx context.Context, gameId uuid.UUID) (dmg.Game, error)
	GetAvailable(ctx context.Context, userId uuid.UUID) ([]dmg.GameMeta, error)
	GetFinished(ctx context.Context, userId uuid.UUID) ([]dmg.GameMeta, error)
	GetLeaderboard(ctx context.Context, limit int) ([]dmg.UserStats, error)
}
