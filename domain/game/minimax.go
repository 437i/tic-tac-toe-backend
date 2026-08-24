package game

import (
	"context"
	"fmt"
)

type Minimax struct{}

func NewMinimax() Minimax {
	return Minimax{}
}

func (ai Minimax) Move(ctx context.Context, game Game) (Game, error) {
	bestI, bestJ, err := getBestMove(ctx, game.Field, game.getHuman())
	if err != nil {
		return Game{}, fmt.Errorf("error getting best move: %w", err)
	}
	game.setCell(bestI, bestJ, getOpponentSign(game.getHuman()))
	game.switchTurn()
	return game, nil
}

func getBestMove(ctx context.Context, field Field, human XO) (int, int, error) {
	bestI, bestJ := -1, -1
	bestScore := -100
	for i := range field {
		for j := range field[i] {
			if field[i][j] == Empty {
				next := field
				next[i][j] = getOpponentSign(human)
				score, err := minimax(ctx, next, 0, false, human)
				if err != nil {
					return 0, 0, err
				}
				if score > bestScore {
					bestScore = score
					bestI, bestJ = i, j
				}
			}
		}
	}
	if bestI == -1 || bestJ == -1 {
		return 0, 0, ErrNoMovesFound
	}
	return bestI, bestJ, nil
}

func minimax(ctx context.Context, field Field, depth int, maximize bool, human XO) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	if winner := getWinner(field); winner != Empty {
		return score(winner, human, depth), nil
	}
	if fieldFull(field) {
		return 0, nil
	}
	if maximize {
		bestScore := -100
		for i := range field {
			for j := range field[i] {
				if field[i][j] == Empty {
					next := field
					next[i][j] = getOpponentSign(human)
					score, err := minimax(ctx, next, depth+1, false, human)
					if err != nil {
						return 0, err
					}
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore, nil
	} else {
		bestScore := 100
		for i := range field {
			for j := range field[i] {
				if field[i][j] == Empty {
					next := field
					next[i][j] = human
					score, err := minimax(ctx, next, depth+1, true, human)
					if err != nil {
						return 0, err
					}
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore, nil
	}
}

func score(winner, human XO, depth int) int {
	switch winner {
	case human:
		return -10 + depth
	case getOpponentSign(human):
		return 10 - depth
	}
	return 0
}
