package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// HealthCheck provides a health check endpoint
// @Summary Health Check
// @Description Check the health status of the application and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Failure 503 {object} map[string]interface{} "Service is unhealthy"
// @Router /health [get]
func HealthCheck(db *pgx.Conn, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"services":  map[string]string{},
		}

		// Check database connection
		if err := db.Ping(ctx); err != nil {
			health["status"] = "unhealthy"
			health["services"].(map[string]string)["database"] = "down"
		} else {
			health["services"].(map[string]string)["database"] = "up"
		}

		// Check Redis connection
		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			health["status"] = "unhealthy"
			health["services"].(map[string]string)["redis"] = "down"
		} else {
			health["services"].(map[string]string)["redis"] = "up"
		}

		// Set appropriate status code
		statusCode := http.StatusOK
		if health["status"] == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(health)
	}
}

// GetStats provides global statistics
// @Summary Get Global Statistics
// @Description Get global statistics including total URLs, clicks, and users
// @Tags statistics
// @Produce json
// @Success 200 {object} map[string]interface{} "Global statistics"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /stats [get]
func GetStats(db *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get total URLs
		totalURLs, err := db.GetTotalURLs(ctx)
		if err != nil {
			http.Error(w, "Failed to get URL statistics", http.StatusInternalServerError)
			return
		}

		// Get total clicks
		totalClicks, err := db.GetTotalClicks(ctx)
		if err != nil {
			http.Error(w, "Failed to get click statistics", http.StatusInternalServerError)
			return
		}

		// Get total users
		totalUsers, err := db.GetTotalUsers(ctx)
		if err != nil {
			http.Error(w, "Failed to get user statistics", http.StatusInternalServerError)
			return
		}

		stats := map[string]interface{}{
			"total_urls":   totalURLs,
			"total_clicks": totalClicks,
			"total_users":  totalUsers,
			"timestamp":    time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
