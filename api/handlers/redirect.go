package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/redis/go-redis/v9"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// RedirectURL redirects to the original URL
// @Summary Redirect to Original URL
// @Description Redirect to the original URL using the short ID and log the click
// @Tags urls
// @Param shortID path string true "Short URL ID"
// @Success 301 "Redirect to original URL"
// @Failure 404 {object} map[string]string "URL not found"
// @Router /{shortID} [get]
func RedirectURL(db *sqlc.Queries, redisClient *redis.Client, geoAPIURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "shortID")
		ctx := r.Context()

		// Check Redis cache
		longURL, err := redisClient.Get(ctx, "url:"+shortID).Result()
		if err == redis.Nil {
			url, err := db.GetURL(ctx, shortID)
			if err != nil {
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}
			longURL = url.LongUrl
			redisClient.Set(ctx, "url:"+shortID, longURL, 24*60*60)
		}

		// Async click logging with background context
		go func() {
			// Use background context to avoid cancellation when request completes
			bgCtx := context.Background()
			db.LogClick(bgCtx, sqlc.LogClickParams{
				ShortID:   pgtype.Text{String: shortID, Valid: true},
				IpAddress: pgtype.Text{String: r.RemoteAddr, Valid: true},
				UserAgent: pgtype.Text{String: r.UserAgent(), Valid: true},
				ClickedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			})
		}()

		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	}
}
