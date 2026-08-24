package game

import (
	com "apg105/web/common"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := com.GetUserIDFromReq(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	var req CreateGameRequest
	if err := com.DecodeJSON(r, &req); err != nil {
		com.WriteError(w, err)
		return
	}
	if err := req.validate(); err != nil {
		com.WriteError(w, err)
		return
	}

	game, err := h.service.CreateGame(r.Context(), req.toDomain(userID))
	if err != nil {
		com.WriteError(w, err)
		return
	}

	createReply := toFull(game)
	com.WriteJSON(w, http.StatusCreated, createReply)
}

func (h *Handler) MakeMove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID, err := com.GetIDFromUrl(params)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	userID, err := com.GetUserIDFromReq(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	var req MoveRequest

	if err := com.DecodeJSON(r, &req); err != nil {
		com.WriteError(w, err)
		return
	}

	game, err := req.toDomain(gameID)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	game, err = h.service.MakeMove(r.Context(), userID, game)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	newGame := toShort(game)
	com.WriteJSON(w, http.StatusOK, newGame)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID, err := com.GetIDFromUrl(params)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	game, err := h.service.GetGame(r.Context(), gameID)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	replyGet := toFull(game)
	com.WriteJSON(w, http.StatusOK, replyGet)
}

func (h *Handler) GetAvailable(w http.ResponseWriter, r *http.Request) {
	userID, err := com.GetUserIDFromReq(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	games, err := h.service.GetAvailable(r.Context(), userID)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	metaListReponse := make([]GameMetaResponse, 0, len(games))
	for _, game := range games {
		metaListReponse = append(metaListReponse, toMeta(game, userID))
	}

	com.WriteJSON(w, http.StatusOK, metaListReponse)
}

func (h *Handler) GetFinished(w http.ResponseWriter, r *http.Request) {
	userID, err := com.GetUserIDFromReq(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	games, err := h.service.GetFinished(r.Context(), userID)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	metaListReponse := make([]GameMetaResponse, 0, len(games))
	for _, game := range games {
		metaListReponse = append(metaListReponse, toMeta(game, userID))
	}

	com.WriteJSON(w, http.StatusOK, metaListReponse)
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	userID, err := com.GetUserIDFromReq(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	params := mux.Vars(r)
	gameID, err := com.GetIDFromUrl(params)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	if err := h.service.Join(r.Context(), gameID, userID); err != nil {
		com.WriteError(w, err)
		return
	}

	com.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "joined",
	})
}

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit, err := com.GetLimitFromURL(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	leaders, err := h.service.GetLeaderboard(r.Context(), limit)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	response := make([]LeaderResponse, 0, len(leaders))
	for _, leader := range leaders {
		response = append(response, toLeader(leader))
	}

	com.WriteJSON(w, http.StatusOK, response)
}
