package user

import (
	"github.com/google/uuid"
)

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type MeResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Login  string    `json:"login"`
}
