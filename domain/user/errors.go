package user

import "errors"

var (
	ErrEmptyLogpass       = errors.New("login/password can't be empty")
	ErrPassTooLong        = errors.New("password too long")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
