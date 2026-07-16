package http

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

type tokenParser interface {
	ParseToken(tokenStr string) (int64, error)
}

func RequireAuth(parser tokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			userID, err := parser.ParseToken(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func userIDFromCtx(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value(userIDKey).(int64)
	return id, ok
}
