package game

import (
	dmg "apg105/domain/game"
	"fmt"

	"github.com/google/uuid"
)

func toRepoField(field dmg.Field) Field {
	result := Field{}
	idx := 0
	for i := range 3 {
		for j := range 3 {
			switch field[i][j] {
			case dmg.X:
				result[idx] = X
			case dmg.O:
				result[idx] = O
			default:
				result[idx] = Empty
			}
			idx++
		}
	}
	return result
}

func (f Field) toDomain() dmg.Field {
	result := dmg.Field{}
	idx := 0
	for i := range 3 {
		for j := range 3 {
			switch f[idx] {
			case X:
				result[i][j] = dmg.X
			case O:
				result[i][j] = dmg.O
			default:
				result[i][j] = dmg.Empty
			}
			idx++
		}
	}
	return result
}

func statusToDomain(status string) (dmg.Status, error) {
	switch status {
	case statusWaitingForOpponent:
		return dmg.WaitingForOpponent, nil
	case statusPlayerXTurn:
		return dmg.PlayerXTurn, nil
	case statusPlayerOTurn:
		return dmg.PlayerOTurn, nil
	case statusPlayerXWon:
		return dmg.PlayerXWon, nil
	case statusPlayerOWon:
		return dmg.PlayerOWon, nil
	case statusDraw:
		return dmg.Draw, nil
	default:
		return 0, fmt.Errorf("unknown game status: %q", status)
	}
}

func modeToDomain(mode string) (dmg.Mode, error) {
	switch mode {
	case modePvP:
		return dmg.PvP, nil
	case modePvE:
		return dmg.PvE, nil
	default:
		return 0, fmt.Errorf("unknown game mode: %q", mode)
	}
}

func uuidPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func uuidValue(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

func toRepoGame(game dmg.Game) Game {
	return Game{
		GameId:  game.GameID,
		Field:   toRepoField(game.Field),
		PlayerX: uuidPtr(game.PlayerX),
		PlayerO: uuidPtr(game.PlayerO),
		Status:  game.Status.String(),
		Mode:    game.Mode.String(),
		Version: game.Version,
	}
}

func (g Game) toDomain() (dmg.Game, error) {
	status, err := statusToDomain(g.Status)
	if err != nil {
		return dmg.Game{}, fmt.Errorf("map status: %w", err)
	}

	mode, err := modeToDomain(g.Mode)
	if err != nil {
		return dmg.Game{}, fmt.Errorf("map mode: %w", err)
	}

	return dmg.Game{
		GameID:    g.GameId,
		Field:     g.Field.toDomain(),
		PlayerX:   uuidValue(g.PlayerX),
		PlayerO:   uuidValue(g.PlayerO),
		Status:    status,
		Mode:      mode,
		CreatedAt: g.CreatedAt,
		Version:   g.Version,
	}, nil
}

func (g GameMeta) toDomain() (dmg.GameMeta, error) {
	status, err := statusToDomain(g.Status)
	if err != nil {
		return dmg.GameMeta{}, fmt.Errorf("map status: %w", err)
	}

	mode, err := modeToDomain(g.Mode)
	if err != nil {
		return dmg.GameMeta{}, fmt.Errorf("map mode: %w", err)
	}

	return dmg.GameMeta{
		GameID:    g.GameID,
		PlayerX:   uuidValue(g.PlayerX),
		PlayerO:   uuidValue(g.PlayerO),
		Status:    status,
		Mode:      mode,
		CreatedAt: g.CreatedAt,
	}, nil
}

func (s UserStats) toDomain() dmg.UserStats {
	return dmg.UserStats{
		UserID:     s.UserID,
		Login:      s.Login,
		TotalGames: s.TotalGames,
		Wins:       s.Wins,
		Losses:     s.Losses,
		Draws:      s.Draws,
		Winrate:    s.Winrate,
	}
}
