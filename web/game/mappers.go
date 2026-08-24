package game

import (
	dmg "apg105/domain/game"
	"apg105/web/common"
	errs "apg105/web/common"
	"fmt"

	"github.com/google/uuid"
)

func (r CreateGameRequest) validate() error {
	if r.Role != X && r.Role != O {
		return errs.ErrInvalidRole
	}
	if r.Mode != PvP && r.Mode != PvE {
		return errs.ErrInvalidMode
	}
	return nil
}

// хреново сделал, можно упростить, но работает
// проверить в другой нейронке
// func (r CreateGameRequest) toDomain(userID uuid.UUID) dmg.GameCreation {
// 	r.PlayerID = userID
// 	game := dmg.GameCreation{}
// 	if r.Mode == PvP {
// 		game.Mode = dmg.PvP
// 		if r.Role == X {
// 			game.PlayerX = r.PlayerID
// 		} else {
// 			game.PlayerO = r.PlayerID
// 		}
// 	} else {
// 		game.Mode = dmg.PvE
// 		if r.Role == X {
// 			game.PlayerX = r.PlayerID
// 			game.PlayerO = uuid.Nil
// 		} else {
// 			game.PlayerX = uuid.Nil
// 			game.PlayerO = r.PlayerID
// 		}
// 	}
// 	return game
// }

func (r CreateGameRequest) toDomain(userID uuid.UUID) dmg.GameCreation {
	game := dmg.GameCreation{}
	if r.Mode == PvP {
		game.Mode = dmg.PvP
	} else {
		game.Mode = dmg.PvE
	}
	if r.Role == X {
		game.PlayerX = userID
	} else {
		game.PlayerO = userID
	}
	return game
}

func (m MoveRequest) toDomain(gameID uuid.UUID) (dmg.Game, error) {
	field, err := m.Field.toDomain()
	if err != nil {
		return dmg.Game{}, fmt.Errorf("mapping field to domain: %w", err)
	}

	return dmg.Game{
		GameID: gameID,
		Field:  field,
	}, nil
}

func toFull(game dmg.Game) FullGameResponse {
	result := FullGameResponse{
		GameID: game.GameID,
		Field:  toWebField(game.Field),
		Status: toWebStatus(
			game.Status,
			game.PlayerX,
			game.PlayerO,
			game.Mode,
		),
		Mode: toWebMode(game.Mode),
	}
	result.PlayerX, result.PlayerO = getPlayersStr(
		game.PlayerX,
		game.PlayerO,
		game.Mode,
	)
	return result
}

func toShort(game dmg.Game) ShortGameResponse {
	result := ShortGameResponse{
		Field: toWebField(game.Field),
		Status: toWebStatus(
			game.Status,
			game.PlayerX,
			game.PlayerO,
			game.Mode,
		),
	}
	result.PlayerX, result.PlayerO = getPlayersStr(
		game.PlayerX,
		game.PlayerO,
		game.Mode,
	)
	return result
}

func toMeta(game dmg.GameMeta, userId uuid.UUID) GameMetaResponse {
	result := GameMetaResponse{
		GameID: game.GameID,
	}
	result.Status = toWebStatus(
		game.Status,
		game.PlayerX,
		game.PlayerO,
		game.Mode,
	)
	result.PlayerX, result.PlayerO = getPlayersStr(
		game.PlayerX,
		game.PlayerO,
		game.Mode,
	)
	if game.PlayerX == userId {
		result.PlayerX = "you"
	}
	if game.PlayerO == userId {
		result.PlayerO = "you"
	}
	result.CreatedAt = game.CreatedAt.Format("2006/01/02 15:04:05")
	return result
}

func getPlayersStr(playerX, playerO uuid.UUID, mode dmg.Mode) (string, string) {
	playerXNew, playerONew := "empty", "empty"
	if playerX == uuid.Nil && mode == dmg.PvE {
		playerXNew = "AI"
	} else if playerX != uuid.Nil {
		playerXNew = playerX.String()
	}
	if playerO == uuid.Nil && mode == dmg.PvE {
		playerONew = "AI"
	} else if playerO != uuid.Nil {
		playerONew = playerO.String()
	}
	return playerXNew, playerONew
}

func toWebField(field dmg.Field) Field {
	res := Field{}
	for i := range field {
		for j := range field[i] {
			if field[i][j] == dmg.Empty {
				res[i][j] = Empty
			}
			if field[i][j] == dmg.X {
				res[i][j] = X
			}
			if field[i][j] == dmg.O {
				res[i][j] = O
			}
		}
	}
	return res
}

func toWebStatus(status dmg.Status, playerX, playerO uuid.UUID, mode dmg.Mode) Status {
	playerXNew, playerONew := getPlayersStr(playerX, playerO, mode)
	switch status {
	case dmg.WaitingForOpponent:
		return WaitingForOpponent
	case dmg.PlayerXTurn:
		return Status(fmt.Sprintf("player %s turn", playerXNew))
	case dmg.PlayerOTurn:
		return Status(fmt.Sprintf("player %s turn", playerONew))
	case dmg.PlayerXWon:
		if playerX == uuid.Nil && mode == dmg.PvE {
			return Status("AI won")
		}
		return Status(fmt.Sprintf("player %s won", playerXNew))
	case dmg.PlayerOWon:
		if playerO == uuid.Nil && mode == dmg.PvE {
			return Status("AI won")
		}
		return Status(fmt.Sprintf("player %s won", playerONew))
	case dmg.Draw:
		return Draw
	}
	return UnknownStatus
}

func toWebMode(mode dmg.Mode) Mode {
	switch mode {
	case dmg.PvP:
		return PvP
	case dmg.PvE:
		return PvE
	}
	return UnknownMode
}

func (f Field) toDomain() (dmg.Field, error) {
	res := dmg.Field{}
	for i := range f {
		for j := range f[i] {
			switch f[i][j] {
			case Empty:
				res[i][j] = dmg.Empty
			case X:
				res[i][j] = dmg.X
			case O:
				res[i][j] = dmg.O
			default:
				return dmg.Field{}, common.ErrInvalidField
			}
		}
	}
	return res, nil
}

func toLeader(user dmg.UserStats) LeaderResponse {
	return LeaderResponse{
		UserID:  user.UserID,
		Login:   user.Login,
		Winrate: user.Winrate,
	}
}
