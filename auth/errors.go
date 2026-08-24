package auth

import "errors"

var (
	ErrEmptyJWTEnvKey           = errors.New("empty jwt env key")
	ErrSecretTooShort           = errors.New("secret must be at least 32 bytes long")
	ErrAccessExpNotPositive     = errors.New("ACCESS_EXP must be positive")
	ErrRefreshExpNotPositive    = errors.New("REFRESH_EXP must be positive")
	ErrRefreshExpLessThanAccess = errors.New("refresh token expiration must be greater than access token expiration")

	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)
