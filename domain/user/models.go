package user

import "github.com/google/uuid"

type SafeUser struct {
	UserID uuid.UUID
	Login  string
}

type User struct {
	UserID   uuid.UUID
	Login    string
	Password string
}

type SignUpRequest struct {
	Login    string
	Password string
}

func NewUser(signup SignUpRequest) User {
	return User{
		UserID:   uuid.New(),
		Login:    signup.Login,
		Password: signup.Password,
	}
}
