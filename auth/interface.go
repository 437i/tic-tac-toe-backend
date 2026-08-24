package auth

import (
	dmu "apg105/domain/user"
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Login(ctx context.Context, login, pass string) (dmu.SafeUser, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (dmu.SafeUser, error)
}

type TokenSigner interface {
	GenerateAccessToken(user dmu.SafeUser) (string, error)
	GenerateRefreshToken(user dmu.SafeUser) (string, error)
}

type TokenParser interface {
	GetIDFromAccessToken(token string) (uuid.UUID, error)
	GetIDFromRefreshToken(token string) (uuid.UUID, error)
}

type TokenManager interface {
	TokenSigner
	TokenParser
}
