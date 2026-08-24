package user

import (
	dmu "apg105/domain/user"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req dmu.SignUpRequest) error
	GetUserLogin(ctx context.Context, userID uuid.UUID) (string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (dmu.SafeUser, error)
}
