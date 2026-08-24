package router

import (
	"apg105/web"
	"apg105/web/auth"
	"apg105/web/game"
	"apg105/web/user"

	"github.com/gorilla/mux"
)

func NewRouter(
	gameH *game.Handler,
	userH *user.Handler,
	authH *auth.Handler,
	authMid AuthMiddleware,
) *mux.Router {
	r := mux.NewRouter()
	r.Use(web.LoggingMiddleware)

	r.HandleFunc("/signup", userH.Signup).Methods("POST").Name("signup")
	r.HandleFunc("/login", authH.Login).Methods("POST").Name("login")
	r.HandleFunc("/refresh/access", authH.AccessRefresh).Methods("POST").Name("access_refresh")
	r.HandleFunc("/refresh/refresh", authH.RefreshRefresh).Methods("POST").Name("refresh_refresh")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(authMid.AuthMiddleware)

	protected.HandleFunc("/leaderboard", gameH.GetLeaderboard).Methods("GET").Name("get_leaders")

	user := protected.PathPrefix("/user").Subrouter()
	user.HandleFunc("/me", userH.Me).Methods("GET").Name("get_me")
	user.HandleFunc("/{id}", userH.GetUser).Methods("GET").Name("get_user")

	game := protected.PathPrefix("/game").Subrouter()
	game.HandleFunc("", gameH.Create).Methods("POST").Name("create_game")
	game.HandleFunc("/available", gameH.GetAvailable).Methods("GET").Name("get_available")
	game.HandleFunc("/history", gameH.GetFinished).Methods("GET").Name("get_history")
	game.HandleFunc("/{id}", gameH.Get).Methods("GET").Name("get_game")
	game.HandleFunc("/{id}/join", gameH.Join).Methods("POST").Name("join")
	game.HandleFunc("/{id}", gameH.MakeMove).Methods("POST").Name("make_move")
	return r
}
