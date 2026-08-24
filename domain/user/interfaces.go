package user

import (
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, userDm User) error
	GetUserByLogin(ctx context.Context, login string) (User, error)
	GetLoginByID(ctx context.Context, userID uuid.UUID) (string, error)
}
