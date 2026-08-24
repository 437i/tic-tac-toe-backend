package auth

import (
	dmu "apg105/domain/user"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtProvider struct {
	accessSecret    []byte
	refreshSecret   []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

func NewJwtProvider() (*JwtProvider, error) {
	access := os.Getenv("ACCESS")
	refresh := os.Getenv("REFRESH")
	if access == "" || refresh == "" {
		return nil, ErrEmptyJWTEnvKey
	}
	if len(access) < 32 || len(refresh) < 32 {
		return nil, ErrSecretTooShort
	}

	accessExpStr := os.Getenv("ACCESS_EXP")
	refreshExpStr := os.Getenv("REFRESH_EXP")

	accessExp := 1 * time.Minute
	refreshExp := 24 * time.Hour

	var err error

	if accessExpStr != "" {
		accessExp, err = time.ParseDuration(accessExpStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ACCESS_EXP: %w", err)
		}
		if accessExp <= 0 {
			return nil, ErrAccessExpNotPositive
		}
	}

	if refreshExpStr != "" {
		refreshExp, err = time.ParseDuration(refreshExpStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REFRESH_EXP: %w", err)
		}
		if refreshExp <= 0 {
			return nil, ErrRefreshExpNotPositive
		}
	}

	if refreshExp <= accessExp {
		return nil, ErrRefreshExpLessThanAccess
	}

	return &JwtProvider{
		accessSecret:    []byte(access),
		refreshSecret:   []byte(refresh),
		accessTokenExp:  accessExp,
		refreshTokenExp: refreshExp,
	}, nil
}

func (p *JwtProvider) GenerateAccessToken(user dmu.SafeUser) (string, error) {
	return p.generateToken(user, p.accessSecret, p.accessTokenExp)
}

func (p *JwtProvider) GenerateRefreshToken(user dmu.SafeUser) (string, error) {
	return p.generateToken(user, p.refreshSecret, p.refreshTokenExp)
}

func (p *JwtProvider) ValidateAccessToken(token string) error {
	_, err := p.parseToken(token, p.accessSecret)
	return err
}

func (p *JwtProvider) ValidateRefreshToken(token string) error {
	_, err := p.parseToken(token, p.refreshSecret)
	return err
}

func (p *JwtProvider) GetIDFromAccessToken(token string) (uuid.UUID, error) {
	return p.getIDFromToken(token, p.accessSecret)
}

func (p *JwtProvider) GetIDFromRefreshToken(token string) (uuid.UUID, error) {
	return p.getIDFromToken(token, p.refreshSecret)
}

func (p *JwtProvider) getIDFromToken(token string, secret []byte) (uuid.UUID, error) {
	claims, err := p.parseToken(token, secret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing token: %w", err)
	}

	return claims.UserID, nil
}

func (p *JwtProvider) generateToken(user dmu.SafeUser, secret []byte, exp time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func (p *JwtProvider) parseToken(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%w: invalid token claims", ErrInvalidToken)
	}
	if claims.UserID == uuid.Nil {
		return nil, fmt.Errorf("%w: user id is nil", ErrInvalidToken)
	}

	return claims, nil
}
