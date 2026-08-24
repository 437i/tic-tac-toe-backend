package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	ai   AI
	repo Repo
}

func NewService(ai AI, repo Repo) *Service {
	return &Service{ai, repo}
}

func (s *Service) CreateGame(ctx context.Context, req GameCreation) (Game, error) {
	game, err := NewGame(req)
	if err != nil {
		return Game{}, fmt.Errorf("creating game: %w", err)
	}

	if game.startsAI() {
		var err error
		game, err = s.ai.Move(ctx, game)
		if err != nil {
			return Game{}, fmt.Errorf("get move: %w", err)
		}
	}

	if err := s.repo.Create(ctx, game); err != nil {
		return Game{}, fmt.Errorf("create game: %w", err)
	}

	return game, nil
}

func (s *Service) MakeMove(ctx context.Context, userID uuid.UUID, game Game) (Game, error) {
	oldGame, err := s.GetGame(ctx, game.GameID)
	if err != nil {
		return Game{}, err
	}

	if oldGame.isOver() {
		return Game{}, ErrGameOver
	}

	if err := game.applyMoveTo(oldGame, userID); err != nil {
		return Game{}, fmt.Errorf("validate game: %w", err)
	}

	if game.isOver() {
		if err := game.finish(); err != nil {
			return Game{}, fmt.Errorf("finishing game: %w", err)
		}
		newVersion, err := s.repo.Update(ctx, game)
		if err != nil {
			return Game{}, fmt.Errorf("save game: %w", err)
		}
		game.Version = newVersion
		return game, nil
	}

	if game.isPvE() {
		aiGame, err := s.ai.Move(ctx, game)
		if err != nil {
			return Game{}, fmt.Errorf("get ai move: %w", err)
		}

		if aiGame.isOver() {
			if err := aiGame.finish(); err != nil {
				return Game{}, fmt.Errorf("finishing ai game: %w", err)
			}
		}
		newVersion, err := s.repo.Update(ctx, aiGame)
		if err != nil {
			return Game{}, fmt.Errorf("save game: %w", err)
		}
		aiGame.Version = newVersion
		return aiGame, nil
	}
	newVersion, err := s.repo.Update(ctx, game)
	if err != nil {
		return Game{}, fmt.Errorf("save game: %w", err)
	}
	game.Version = newVersion
	return game, nil
}

func (s *Service) Join(ctx context.Context, gameID, userID uuid.UUID) error {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game: %w", err)
	}

	if err := game.validateJoin(userID); err != nil {
		return fmt.Errorf("validate join: %w", err)
	}

	if err := game.addPlayer(userID); err != nil {
		return fmt.Errorf("add player to game: %w", err)
	}
	newVersion, err := s.repo.Update(ctx, game)
	if err != nil {
		return fmt.Errorf("save game: %w", err)
	}
	game.Version = newVersion
	return nil
}

func (s *Service) GetGame(ctx context.Context, gameId uuid.UUID) (Game, error) {
	return s.repo.Get(ctx, gameId)
}

func (s *Service) GetAvailable(ctx context.Context, userId uuid.UUID) ([]GameMeta, error) {
	return s.repo.GetAvailable(ctx, userId)
}

func (s *Service) GetFinished(ctx context.Context, userId uuid.UUID) ([]GameMeta, error) {
	return s.repo.GetFinished(ctx, userId)
}

func (s *Service) GetLeaderboard(ctx context.Context, limit int) ([]UserStats, error) {
	const (
		minLimit = 1
		maxLimit = 100
	)

	limit = min(max(limit, minLimit), maxLimit)

	return s.repo.GetLeaderboard(ctx, limit)
}
