package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

func AuthMiddleware(db *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKeyStr := r.Header.Get("X-API-Key")
			apiKey, err := uuid.Parse(apiKeyStr)
			if err != nil {
				http.Error(w, "Invalid API key format", http.StatusUnauthorized)
				return
			}

			key, err := db.GetAPIKey(r.Context(), apiKey)
			if err != nil {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", key.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
