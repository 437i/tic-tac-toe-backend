package auth

import (
	"apg105/auth"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req auth.JWTRequest) (auth.JWTResponse, error)
	RefreshAccessToken(ctx context.Context, refresh string) (auth.JWTResponse, error)
	RefreshRefreshToken(ctx context.Context, refresh string) (auth.JWTResponse, error)
}
