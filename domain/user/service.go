package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo}
}

func (s *Service) Register(ctx context.Context, req SignUpRequest) error {
	if req.Login == "" || req.Password == "" {
		return ErrEmptyLogpass
	}

	if len(req.Password) > 72 {
		return ErrPassTooLong
	}

	user := NewUser(req)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	user.Password = string(hashedPassword)

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, login, pass string) (SafeUser, error) {
	if login == "" || pass == "" {
		return SafeUser{}, ErrEmptyLogpass
	}

	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return SafeUser{}, ErrInvalidCredentials
		}
		return SafeUser{}, fmt.Errorf("get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	switch {
	case err == nil:
		return SafeUser{
			UserID: user.UserID,
			Login:  user.Login,
		}, nil
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return SafeUser{}, ErrInvalidCredentials
	default:
		return SafeUser{}, fmt.Errorf("compare password hash: %w", err)
	}
}

func (s *Service) GetUserLogin(ctx context.Context, userID uuid.UUID) (string, error) {
	return s.repo.GetLoginByID(ctx, userID)
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (SafeUser, error) {
	login, err := s.repo.GetLoginByID(ctx, userID)
	if err != nil {
		return SafeUser{}, fmt.Errorf("getting login by id: %w", err)
	}

	return SafeUser{
		UserID: userID,
		Login:  login,
	}, nil
}
