package auth

import (
	"apg105/web/common"
	"net/http"
)

type Handler struct {
	service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var jwtRequest JWTRequest

	if err := common.DecodeJSON(r, &jwtRequest); err != nil {
		common.WriteError(w, err)
		return
	}

	response, err := h.service.Login(r.Context(), jwtRequest.toDomain())
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, toWebResponse(response))
}

func (h *Handler) AccessRefresh(w http.ResponseWriter, r *http.Request) {
	var refreshReq RefreshJWTRequest

	if err := common.DecodeJSON(r, &refreshReq); err != nil {
		common.WriteError(w, err)
		return
	}

	response, err := h.service.RefreshAccessToken(r.Context(), refreshReq.RefreshToken)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, toWebResponse(response))
}

func (h *Handler) RefreshRefresh(w http.ResponseWriter, r *http.Request) {
	var refreshReq RefreshJWTRequest

	if err := common.DecodeJSON(r, &refreshReq); err != nil {
		common.WriteError(w, err)
		return
	}

	response, err := h.service.RefreshRefreshToken(r.Context(), refreshReq.RefreshToken)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, toWebResponse(response))
}
