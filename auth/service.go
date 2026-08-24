package auth

import (
	dmu "apg105/domain/user"
	"context"
	"fmt"
)

type AuthService struct {
	jwt         TokenManager
	userService UserService
}

func NewAuthService(jwt TokenManager, userService UserService) *AuthService {
	return &AuthService{jwt, userService}
}

func (s *AuthService) Login(ctx context.Context, req JWTRequest) (JWTResponse, error) {
	user, err := s.userService.Login(ctx, req.Login, req.Password)
	if err != nil {
		return JWTResponse{}, fmt.Errorf("user service login: %w", err)
	}
	access, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return JWTResponse{}, fmt.Errorf("generate access token: %w", err)
	}
	refresh, err := s.jwt.GenerateRefreshToken(user)
	if err != nil {
		return JWTResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}
	return JWTResponse{
		Type:         "Bearer",
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refresh string) (JWTResponse, error) {
	return s.refresh(ctx, refresh, false)
}

func (s *AuthService) RefreshRefreshToken(ctx context.Context, refresh string) (JWTResponse, error) {
	return s.refresh(ctx, refresh, true)
}

func (s *AuthService) GetUserFromToken(ctx context.Context, accessToken string) (dmu.SafeUser, error) {
	id, err := s.jwt.GetIDFromAccessToken(accessToken)
	if err != nil {
		return dmu.SafeUser{}, fmt.Errorf("get id from access token: %w", err)
	}
	user, err := s.userService.GetUserByID(ctx, id)
	if err != nil {
		return dmu.SafeUser{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *AuthService) refresh(ctx context.Context, refreshToken string, rotate bool) (JWTResponse, error) {
	id, err := s.jwt.GetIDFromRefreshToken(refreshToken)
	if err != nil {
		return JWTResponse{}, fmt.Errorf("get id from refresh token: %w", err)
	}
	user, err := s.userService.GetUserByID(ctx, id)
	if err != nil {
		return JWTResponse{}, fmt.Errorf("get user: %w", err)
	}
	accessToken, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return JWTResponse{}, fmt.Errorf("generate access token: %w", err)
	}
	if rotate {
		refreshToken, err = s.jwt.GenerateRefreshToken(user)
		if err != nil {
			return JWTResponse{}, fmt.Errorf("generate refresh token: %w", err)
		}
	}
	return JWTResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
