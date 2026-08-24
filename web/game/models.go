package game

import (
	"github.com/google/uuid"
)

type XO string

const (
	X     XO = "X"
	O     XO = "O"
	Empty XO = " "
)

type Status string

const (
	WaitingForOpponent Status = "waiting for opponent"
	Draw               Status = "draw"
	UnknownStatus      Status = "unknown status"
)

type Mode string

const (
	PvE         Mode = "PvE"
	PvP         Mode = "PvP"
	UnknownMode Mode = "unknown mode"
)

type Field [3][3]XO

type CreateGameRequest struct {
	Role XO   `json:"role"`
	Mode Mode `json:"mode"`
}

type MoveRequest struct {
	Field Field `json:"field"`
}

type FullGameResponse struct {
	GameID  uuid.UUID `json:"game_id"`
	Field   Field     `json:"field"`
	PlayerX string    `json:"player_x"`
	PlayerO string    `json:"player_o"`
	Status  Status    `json:"status"`
	Mode    Mode      `json:"mode"`
}

type ShortGameResponse struct {
	Field   Field  `json:"field"`
	PlayerX string `json:"player_x"`
	PlayerO string `json:"player_o"`
	Status  Status `json:"status"`
}

type GameMetaResponse struct {
	GameID    uuid.UUID `json:"game_id"`
	Status    Status    `json:"status"`
	PlayerX   string    `json:"player_x"`
	PlayerO   string    `json:"player_o"`
	CreatedAt string    `json:"created_at"`
}

type LeaderResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Login   string    `json:"login"`
	Winrate float64   `json:"winrate"`
}
