package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// AnalyticsResponse represents the analytics response
type AnalyticsResponse map[string]int

// GetAnalytics gets analytics for a specific URL
// @Summary Get URL Analytics
// @Description Get click analytics for a specific URL owned by the authenticated user
// @Tags analytics
// @Security ApiKeyAuth
// @Param shortID path string true "Short URL ID"
// @Produce json
// @Success 200 {object} AnalyticsResponse "Analytics data by location"
// @Failure 401 {object} map[string]string "Unauthorized or URL not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/analytics/{shortID} [get]
func GetAnalytics(db *sqlc.Queries, geoAPIURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "shortID")
		userIDStr, _ := r.Context().Value("user_id").(string)

		// Parse userID from context
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Ensure user owns the URL
		url, err := db.GetURL(r.Context(), shortID)
		if err != nil || (url.UserID.Valid && url.UserID.Bytes != userID) {
			http.Error(w, "Unauthorized or URL not found", http.StatusUnauthorized)
			return
		}

		clicks, _ := db.ListClicks(r.Context(), pgtype.Text{String: shortID, Valid: true})
		analytics := map[string]int{}
		for _, click := range clicks {
			location := getGeoLocation(click.IpAddress.String, geoAPIURL)
			analytics[location]++
		}
		json.NewEncoder(w).Encode(analytics)
	}
}

func getGeoLocation(ip, geoAPIURL string) string {
	if ip == "" {
		return "unknown"
	}

	// Create request to the geo API
	resp, err := http.Get(geoAPIURL + "/" + ip)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	// Parse the response
	var result struct {
		Country string `json:"country"`
		City    string `json:"city"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "unknown"
	}

	// Check if the request was successful
	if result.Status != "success" {
		return "unknown"
	}

	// Return country or city+country if available
	if result.City != "" {
		return result.City + ", " + result.Country
	}
	if result.Country != "" {
		return result.Country
	}
	return "unknown"
}
