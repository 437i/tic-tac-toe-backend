package game

import (
	"fmt"

	"github.com/google/uuid"
)

func (g Game) getHuman() XO {
	if g.PlayerX == uuid.Nil {
		return O
	}
	return X
}

func getOpponentSign(player XO) XO {
	if player == X {
		return O
	}
	return X
}

func (g Game) isOver() bool {
	return getWinner(g.Field) != Empty || fieldFull(g.Field)
}

func (g Game) compareFieldWith(old Game, player XO) error {
	diffCount := 0
	var err error
	for i := range 3 {
		for j := range 3 {
			diffCount, err = validateCell(
				old.Field[i][j],
				g.Field[i][j],
				player,
				diffCount)
			if err != nil {
				return fmt.Errorf("validate cell: %w", err)
			}
		}
	}
	if diffCount > 1 {
		return ErrManyChanges
	} else if diffCount == 0 {
		return ErrNoChanges
	}
	return nil
}

func validateCell(old, new, player XO, counter int) (int, error) {
	if old != new {
		if old != Empty {
			return 0, ErrCheating
		}
		if new != player {
			return 0, ErrWrongSign
		}
		counter++
	}
	return counter, nil
}

func fieldFull(field Field) bool {
	for i := range field {
		for j := range field[i] {
			if field[i][j] == Empty {
				return false
			}
		}
	}
	return true
}

func getWinner(field Field) XO {
	for i := range field {
		if field[i][0] != Empty &&
			field[i][0] == field[i][1] &&
			field[i][1] == field[i][2] {
			return field[i][0]
		}

		if field[0][i] != Empty &&
			field[0][i] == field[1][i] &&
			field[1][i] == field[2][i] {
			return field[0][i]
		}
	}

	if field[0][0] != Empty &&
		field[0][0] == field[1][1] &&
		field[1][1] == field[2][2] {
		return field[0][0]
	}

	if field[0][2] != Empty &&
		field[0][2] == field[1][1] &&
		field[1][1] == field[2][0] {
		return field[0][2]
	}

	return Empty
}

func (g *Game) syncFrom(oldGame Game) {
	g.PlayerX = oldGame.PlayerX
	g.PlayerO = oldGame.PlayerO
	g.Status = oldGame.Status
	g.Mode = oldGame.Mode
	g.Version = oldGame.Version
}

func getEndGameStatus(winner XO) (Status, error) {
	switch winner {
	case X:
		return PlayerXWon, nil
	case O:
		return PlayerOWon, nil
	case Empty:
		return Draw, nil
	}
	return Status(-1), ErrUnknownStatusFromWinner
}

func (g Game) validateJoin(userID uuid.UUID) error {
	if g.Mode != PvP ||
		g.Status != WaitingForOpponent ||
		(g.PlayerX != uuid.Nil && g.PlayerO != uuid.Nil) ||
		g.PlayerX == userID ||
		g.PlayerO == userID {
		return ErrGameUnavailableToJoin
	}
	return nil
}

func (g *Game) addPlayer(userID uuid.UUID) error {
	switch {
	case g.PlayerO == uuid.Nil:
		g.PlayerO = userID
	case g.PlayerX == uuid.Nil:
		g.PlayerX = userID
	default:
		return ErrGameUnavailableToJoin
	}
	g.Status = PlayerXTurn
	return nil
}

func (g Game) startsAI() bool {
	return g.Mode == PvE && g.getHuman() == O
}

func (g *Game) switchTurn() {
	if g.Status == PlayerOTurn {
		g.Status = PlayerXTurn
		return
	}
	g.Status = PlayerOTurn
}

func (g Game) isPvE() bool {
	return g.Mode == PvE
}

func (g *Game) applyMoveTo(oldGame Game, playerID uuid.UUID) error {
	if err := g.validateWith(oldGame, playerID); err != nil {
		return fmt.Errorf("validate game: %w", err)
	}
	g.syncFrom(oldGame)
	g.switchTurn()
	return nil
}

func (g Game) validateWith(oldGame Game, playerID uuid.UUID) error {
	if oldGame.Status == WaitingForOpponent {
		return ErrGameNotStarted
	}

	if oldGame.Status != PlayerXTurn &&
		oldGame.Status != PlayerOTurn {
		return ErrGameOver
	}

	if oldGame.wrongTurnFor(playerID) {
		return ErrWrongSign
	}

	if err := g.compareFieldWith(oldGame, oldGame.nowTurn()); err != nil {
		return fmt.Errorf("compare field: %w", err)
	}
	return nil
}

func (g Game) wrongTurnFor(playerID uuid.UUID) bool {
	return g.Status == PlayerOTurn && g.PlayerO != playerID ||
		g.Status == PlayerXTurn && g.PlayerX != playerID
}

func (g Game) nowTurn() XO {
	switch g.Status {
	case PlayerXTurn:
		return X
	default:
		return O
	}
}

func (g *Game) finish() error {
	status, err := getEndGameStatus(getWinner(g.Field))
	if err != nil {
		return fmt.Errorf("get endgame status: %w", err)
	}
	g.Status = status
	return nil
}

func (g *Game) setCell(i, j int, sign XO) {
	g.Field[i][j] = sign
}
