package user

import "github.com/google/uuid"

const usersTable = "game_users"

const (
	colUserID = "user_id"
	colLogin = "login"
	colPassword = "password"
)

type User struct {
	UserID   uuid.UUID `db:"user_id"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}
