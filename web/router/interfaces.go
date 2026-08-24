package router

import "net/http"

type AuthMiddleware interface {
	AuthMiddleware(next http.Handler) http.Handler
}
