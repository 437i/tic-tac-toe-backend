package game

import "errors"

var (
	ErrBothPlayersNil          = errors.New("both players nil")
	ErrBothPlayersNotNil       = errors.New("both players non-nil")
	ErrUnknownMode             = errors.New("unknown game mode")
	ErrUnknownStatusFromWinner = errors.New("unknown status from winner")

	ErrGameOver              = errors.New("game already over")
	ErrManyChanges           = errors.New("too many changes")
	ErrNoChanges             = errors.New("no changes")
	ErrCheating              = errors.New("cheating")
	ErrWrongSign             = errors.New("expected another sign")
	ErrNoMovesFound          = errors.New("no move found")
	ErrGameUnavailableToJoin = errors.New("game unavailable to join")
	ErrGameNotStarted        = errors.New("game was not started")

	ErrGameNotFound           = errors.New("game not found")
	ErrConcurrentModification = errors.New("game was modified concurrently")
)
