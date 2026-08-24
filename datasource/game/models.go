package game

import (
	"time"

	"github.com/google/uuid"
)

const gamesTable = "games"

const (
	colGameID    = "game_id"
	colField     = "field"
	colPlayerX   = "player_x"
	colPlayerO   = "player_o"
	colStatus    = "status"
	colMode      = "mode"
	colCreatedAt = "created_at"
	colVersion   = "version"
)

const (
	statusWaitingForOpponent = "WaitingForOpponent"
	statusPlayerXTurn        = "PlayerXTurn"
	statusPlayerOTurn        = "PlayerOTurn"
	statusPlayerXWon         = "PlayerXWon"
	statusPlayerOWon         = "PlayerOWon"
	statusDraw               = "Draw"
)

const (
	modePvP = "PvP"
	modePvE = "PvE"
)

type XO int

const (
	Empty XO = iota
	X
	O
)

type Field [9]XO

type Game struct {
	GameId    uuid.UUID  `db:"game_id"`
	Field     Field      `db:"field"`
	PlayerX   *uuid.UUID `db:"player_x"`
	PlayerO   *uuid.UUID `db:"player_o"`
	Status    string     `db:"status"`
	Mode      string     `db:"mode"`
	CreatedAt time.Time  `db:"created_at"`
	Version   int        `db:"version"`
}

type GameMeta struct {
	GameID    uuid.UUID  `db:"game_id"`
	PlayerX   *uuid.UUID `db:"player_x"`
	PlayerO   *uuid.UUID `db:"player_o"`
	Status    string     `db:"status"`
	Mode      string     `db:"mode"`
	CreatedAt time.Time  `db:"created_at"`
}

type UserStats struct {
	UserID     uuid.UUID `db:"user_id"`
	Login      string    `db:"login"`
	TotalGames int       `db:"total_games"`
	Wins       int       `db:"wins"`
	Losses     int       `db:"losses"`
	Draws      int       `db:"draws"`
	Winrate    float64   `db:"winrate"`
}
