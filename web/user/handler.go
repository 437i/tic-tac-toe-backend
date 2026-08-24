package user

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

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := com.DecodeJSON(r, &req); err != nil {
		com.WriteError(w, err)
		return
	}

	err := h.service.Register(r.Context(), req.toDomain())
	if err != nil {
		com.WriteError(w, err)
		return
	}

	com.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "success",
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := com.GetIDFromUrl(params)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	login, err := h.service.GetUserLogin(r.Context(), id)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	com.WriteJSON(w, http.StatusOK, map[string]string{
		"login": login,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := com.GetUserIDFromReq(r)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		com.WriteError(w, err)
		return
	}

	com.WriteJSON(w, http.StatusOK, toMe(user))
}
