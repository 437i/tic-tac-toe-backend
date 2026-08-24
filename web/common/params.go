package common

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type userIDKey struct{}

func GetIDFromUrl(params map[string]string) (uuid.UUID, error) {
	idStr, ok := params["id"]
	if !ok {
		return uuid.Nil, ErrNoIdParam
	}

	id, err := uuid.Parse(idStr)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, ErrInvalidId
	}

	return id, nil
}

func GetUserIDFromReq(r *http.Request) (uuid.UUID, error) {
	userIDVal := r.Context().Value(userIDKey{})
	if userIDVal == nil {
		return uuid.Nil, ErrUserIDNotFoundInCtx
	}

	id, ok := userIDVal.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidId
	}

	return id, nil
}

func SetUserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func GetLimitFromURL(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get("limit")

	limit := 10

	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, ErrInvalidLimit
		}
	}

	return limit, nil
}
