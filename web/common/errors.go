package common

import (
	"apg105/auth"
	"apg105/domain/game"
	"apg105/domain/user"
	"errors"
	"log"
	"net/http"
)

var (
	ErrInvalidField = errors.New("invalid field")
	ErrNoIdParam    = errors.New("no id param")
	ErrInvalidId    = errors.New("invalid id")

	ErrUserIDNotFoundInCtx = errors.New("User ID not found in context")

	ErrInvalidRole  = errors.New("invalid role")
	ErrInvalidMode  = errors.New("invalid mode")
	ErrInvalidLimit = errors.New("invalid limit")

	ErrInvalidJSON = errors.New("invalid json")
)

func MapError(err error) (int, string) {
	switch {
	// game
	case errors.Is(err, game.ErrGameNotFound):
		return http.StatusNotFound, "game not found"
	// user
	case errors.Is(err, user.ErrUserAlreadyExists):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, user.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	// game
	case errors.Is(err, game.ErrGameNotStarted):
		return http.StatusConflict, "game has not started"
	case errors.Is(err, game.ErrGameOver):
		return http.StatusBadRequest, "game already over"
	case errors.Is(err, game.ErrManyChanges):
		return http.StatusBadRequest, "invalid field"
	case errors.Is(err, game.ErrNoChanges):
		return http.StatusBadRequest, "invalid field"
	case errors.Is(err, game.ErrCheating):
		return http.StatusBadRequest, "cheating!"
	case errors.Is(err, game.ErrWrongSign):
		return http.StatusBadRequest, "expected another sign"
	case errors.Is(err, game.ErrGameUnavailableToJoin):
		return http.StatusConflict, "game unavailable to join"
	case errors.Is(err, game.ErrConcurrentModification):
		return http.StatusConflict, "game was modified, retry"
	// user
	case errors.Is(err, user.ErrEmptyLogpass):
		return http.StatusBadRequest, "empty login/password"
	case errors.Is(err, user.ErrPassTooLong):
		return http.StatusBadRequest, "password too long"
	case errors.Is(err, user.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid credentials"
	// auth
	case errors.Is(err, auth.ErrInvalidToken):
		return http.StatusUnauthorized, "invalid token"
	case errors.Is(err, auth.ErrTokenExpired):
		return http.StatusUnauthorized, "token expired"
	// web
	case errors.Is(err, ErrInvalidField):
		return http.StatusBadRequest, "invalid field"
	case errors.Is(err, ErrNoIdParam):
		return http.StatusBadRequest, "no id param in url"
	case errors.Is(err, ErrInvalidId):
		return http.StatusBadRequest, "invalid id"
	case errors.Is(err, ErrUserIDNotFoundInCtx):
		return http.StatusUnauthorized, "unauthorized"
	case errors.Is(err, ErrInvalidRole):
		return http.StatusBadRequest, "invalid role"
	case errors.Is(err, ErrInvalidMode):
		return http.StatusBadRequest, "invalid mode"
	case errors.Is(err, ErrInvalidLimit):
		return http.StatusBadRequest, "invalid limit"
	case errors.Is(err, ErrInvalidJSON):
		return http.StatusBadRequest, "invalid json"

	default:
		log.Println(err.Error())
		return http.StatusInternalServerError, "internal server error"
	}
}
