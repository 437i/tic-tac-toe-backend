package datasource

// import (
// 	dm "apg104/domain"
// 	"context"
// 	"sync"

// 	"github.com/google/uuid"
// )

// var ModuleMap = fx.Module("xoMap",
// 	fx.Provide(ds.NewStorage),
// 	fx.Provide(
// 		fx.Annotate(
// 			ds.NewMapGameRepo,
// 			fx.As(new(dom.GameRepo)),
// 		),
// 	),
// 	CommonModule,
// )

// type gameRepoMapImpl struct {
// 	repo *sync.Map
// }

// func NewStorage() *sync.Map {
// 	return &sync.Map{}
// }

// func NewMapGameRepo(storage *sync.Map) GameRepo {
// 	return &gameRepoMapImpl{storage}
// }

// func (r *gameRepoMapImpl) Create(ctx context.Context, game dm.Game) error {
// 	if _, ok := r.repo.LoadOrStore(game.ID, gameToRepo(game)); ok {
// 		return ErrGameAlreadyExists
// 	}
// 	return nil
// }

// func (r *gameRepoMapImpl) Exists(ctx context.Context, id uuid.UUID) bool {
// 	_, ok := r.repo.Load(id)
// 	return ok
// }

// func (r *gameRepoMapImpl) Get(ctx context.Context, id uuid.UUID) (dm.Game, error) {
// 	val, ok := r.repo.Load(id)
// 	if !ok {
// 		return dm.Game{}, ErrGameNotFound
// 	}
// 	game, ok := val.(Game)
// 	if !ok {
// 		return dm.Game{}, ErrGameNotFound
// 	}
// 	return game.toDomain(), nil
// }

// func (r *gameRepoMapImpl) Update(ctx context.Context, game dm.Game) error {
// 	if !r.Exists(ctx, game.ID) {
// 		return ErrGameNotFound
// 	}
// 	r.repo.Store(game.ID, gameToRepo(game))
// 	return nil
// }

// func (r *gameRepoPGImpl) GetSuper(ctx context.Context, userId uuid.UUID) ([]dm.Game, error) {
// 	// исключаем законченные игры
// 	notFinishedFilter := sq.NotEq{"status": []string{"PlayerXWon", "PlayerOWon", "Draw"}}
// 	// выбираем игры с ии где участвует пользователь
// 	playingVsAiFilter := sq.And{
// 		sq.Eq{"mode": "PvE"},
// 		sq.Or{
// 			sq.Eq{"player_x": userId},
// 			sq.Eq{"player_o": userId},
// 		},
// 	}
// 	// выбираем игры с человеком где участвует пользователь
// 	playingVsUserFilter := sq.And{
// 		sq.Eq{"mode": "PvP"},
// 		sq.Or{
// 			sq.Eq{"player_x": userId},
// 			sq.Eq{"player_o": userId},
// 		},
// 	}
// 	// выбираем игры к которым может присоединиться пользователь
// 	waitingForOpponentsFilter := sq.And{
// 		sq.Eq{"mode": "PvP"},
// 		sq.Eq{"status": "WaitingForOpponent"},
// 		sq.NotEq{"player_x": userId},
// 		sq.NotEq{"player_o": userId},
// 	}
// 	sql, args, err := r.queryBuilder.
// 		Select("game_id", "player_x", "player_o", "status", "mode").
// 		From(gamesTable).
// 		Where(notFinishedFilter).
// 		Where(sq.Or{
// 			playingVsAiFilter,
// 			playingVsUserFilter,
// 			waitingForOpponentsFilter,
// 		}).
// 		ToSql()
// 	if err != nil {
// 		return nil, fmt.Errorf("building get query: %w", err)
// 	}
// 	rows, err := r.pool.Query(ctx, sql, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("executing get query: %w", err)
// 	}
// 	defer rows.Close()
// 	games := make([]dm.Game, 0)
// 	for rows.Next() {
// 		game := GameSuperMeta{}
// 		if err := rows.Scan(
// 			&game.GameID,
// 			&game.PlayerX,
// 			&game.PlayerO,
// 			&game.Status,
// 			&game.Mode,
// 		); err != nil {
// 			return nil, fmt.Errorf("scanning row: %w", err)
// 		}
// 		games = append(games, game.toDomain())
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterating rows: %w", err)
// 	}
// 	return games, nil
// }
// func (s *gameServiceImpl) GetSuper(ctx context.Context, userId uuid.UUID) ([]Game, error) {
// 	return s.repo.GetSuper(ctx, userId)
// }
// func toSuperMeta(game dm.Game, userId uuid.UUID) GameSuperMetaResponse {
// 	result := GameSuperMetaResponse{
// 		GameID: game.ID,
// 		Status: toWebStatus(game),
// 		Mode:   toWebMode(game.Mode),
// 	}
// 	result.PlayerX, result.PlayerO = getPlayersStr(game)
// 	if game.PlayerX == userId {
// 		result.PlayerX = "you"
// 	}
// 	if game.PlayerO == userId {
// 		result.PlayerO = "you"
// 	}
// 	return result
// }

// func (h *GameHandler) GetSuper(w http.ResponseWriter, r *http.Request) {
// 	// получаем id пользователя
// 	userId, err := getUserIdFromReq(r)
// 	if err != nil {
// 		status, msg := mapErrorToHTTP(err)
// 		writeError(w, status, msg)
// 		return
// 	}
// 	// получаем ВСЕ доступные пользователю игры
// 	games, err := h.service.GetSuper(r.Context(), userId)
// 	if err != nil {
// 		status, msg := mapErrorToHTTP(err)
// 		writeError(w, status, msg)
// 		return
// 	}
// 	// формируем ответ
// 	metaListReponse := make([]GameSuperMetaResponse, 0)
// 	for _, game := range games {
// 		metaListReponse = append(metaListReponse, toSuperMeta(game, userId))
// 	}
// 	// отправляем
// 	writeJSON(w, http.StatusOK, metaListReponse)
// }
