//go:generate stringer -type=Status,Mode
package game

import (
	"time"

	"github.com/google/uuid"
)

type XO int

const (
	Empty XO = iota
	X
	O
)

type Field [3][3]XO

type Mode int

const (
	PvP Mode = 1
	PvE Mode = 2
)

type Status int

const (
	WaitingForOpponent Status = iota
	PlayerXTurn
	PlayerOTurn
	PlayerXWon
	PlayerOWon
	Draw
)

type GameCreation struct {
	PlayerX uuid.UUID
	PlayerO uuid.UUID
	Mode    Mode
}

type Game struct {
	GameID    uuid.UUID
	Field     Field
	PlayerX   uuid.UUID
	PlayerO   uuid.UUID
	Status    Status
	Mode      Mode
	CreatedAt time.Time
	Version   int
}

func NewGame(req GameCreation) (Game, error) {
	if req.PlayerO == uuid.Nil && req.PlayerX == uuid.Nil {
		return Game{}, ErrBothPlayersNil
	}
	if req.PlayerO != uuid.Nil && req.PlayerX != uuid.Nil {
		return Game{}, ErrBothPlayersNotNil
	}
	if req.Mode != PvP && req.Mode != PvE {
		return Game{}, ErrUnknownMode
	}
	game := Game{}
	game.GameID = uuid.New()
	game.Field = Field{}
	game.PlayerX = req.PlayerX
	game.PlayerO = req.PlayerO
	if req.Mode == PvP {
		game.Status = WaitingForOpponent
	} else {
		game.Status = PlayerXTurn
	}
	game.Mode = req.Mode
	game.Version = 1
	return game, nil
}

type GameMeta struct {
	GameID    uuid.UUID
	PlayerX   uuid.UUID
	PlayerO   uuid.UUID
	Status    Status
	Mode      Mode
	CreatedAt time.Time
}

type UserStats struct {
	UserID     uuid.UUID
	Login      string
	TotalGames int
	Wins       int
	Losses     int
	Draws      int
	Winrate    float64
}
