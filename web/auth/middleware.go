package auth

import (
	"apg105/auth"
	"apg105/web/common"
	"net/http"
	"strings"
)

type Authenticator struct {
	jwt auth.TokenParser
}

func NewAuthenticator(jwt auth.TokenParser) *Authenticator {
	return &Authenticator{jwt}
}

func (a *Authenticator) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.WriteCustomError(w, http.StatusUnauthorized, "missing or invalid auth header")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			common.WriteCustomError(w, http.StatusUnauthorized, "missing or invalid auth header")
			return
		}
		accessToken := parts[1]
		if accessToken == "" {
			common.WriteCustomError(w, http.StatusUnauthorized, "missing or invalid access token")
			return
		}
		userID, err := a.jwt.GetIDFromAccessToken(accessToken)
		if err != nil {
			common.WriteError(w, err)
			return
		}
		ctx := common.SetUserIDToContext(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
