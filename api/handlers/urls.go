package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// URLInfo represents URL information
type URLInfo struct {
	ShortID    string     `json:"short_id" example:"abc123"`
	LongURL    string     `json:"long_url" example:"https://example.com"`
	CreatedAt  time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	ClickLimit *int32     `json:"click_limit,omitempty" example:"100"`
}

// ListURLsResponse represents the response for listing URLs
type ListURLsResponse struct {
	URLs []URLInfo `json:"urls"`
}

// UpdateURLRequest represents the request body for updating a URL
type UpdateURLRequest struct {
	LongURL    *string    `json:"long_url,omitempty" example:"https://new-example.com"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	ClickLimit *int       `json:"click_limit,omitempty" example:"200"`
}

// ListUserURLs lists all URLs for the authenticated user
// @Summary List User URLs
// @Description List all URLs created by the authenticated user
// @Tags urls
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} ListURLsResponse "URLs retrieved successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/urls [get]
func ListUserURLs(db *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		urls, err := db.ListUserURLs(r.Context(), sqlc.UUIDToNullable(&userID))
		if err != nil {
			http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
			return
		}

		// Convert to a more user-friendly format
		var response []map[string]interface{}
		for _, url := range urls {
			urlData := map[string]interface{}{
				"short_id":   url.ShortID,
				"long_url":   url.LongUrl,
				"created_at": url.CreatedAt.Time,
			}

			if url.ExpiresAt.Valid {
				urlData["expires_at"] = url.ExpiresAt.Time
			}

			if url.ClickLimit.Valid {
				urlData["click_limit"] = url.ClickLimit.Int32
			}

			response = append(response, urlData)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"urls": response,
		})
	}
}

// DeleteURL deletes a URL for the authenticated user
// @Summary Delete URL
// @Description Delete a URL owned by the authenticated user
// @Tags urls
// @Security ApiKeyAuth
// @Param shortID path string true "Short URL ID"
// @Produce json
// @Success 200 {object} map[string]string "URL deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to delete URL or URL not found"
// @Router /api/urls/{shortID} [delete]
func DeleteURL(db *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "shortID")
		if shortID == "" {
			http.Error(w, "Short ID is required", http.StatusBadRequest)
			return
		}

		userIDStr, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = db.DeleteURL(r.Context(), sqlc.DeleteURLParams{
			ShortID: shortID,
			UserID:  sqlc.UUIDToNullable(&userID),
		})

		if err != nil {
			http.Error(w, "Failed to delete URL or URL not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "URL deleted successfully",
		})
	}
}

// UpdateURL updates a URL for the authenticated user
// @Summary Update URL
// @Description Update URL settings for a URL owned by the authenticated user
// @Tags urls
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param shortID path string true "Short URL ID"
// @Param url body UpdateURLRequest true "URL update information"
// @Success 200 {object} URLInfo "URL updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "URL not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/urls/{shortID} [put]
func UpdateURL(db *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "shortID")
		if shortID == "" {
			http.Error(w, "Short ID is required", http.StatusBadRequest)
			return
		}

		userIDStr, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var input struct {
			LongURL    *string    `json:"long_url,omitempty"`
			ExpiresAt  *time.Time `json:"expires_at,omitempty"`
			ClickLimit *int       `json:"click_limit,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get current URL to use as defaults for unspecified fields
		currentURL, err := db.GetURL(r.Context(), shortID)
		if err != nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		// Check if user owns the URL
		if !currentURL.UserID.Valid || currentURL.UserID.Bytes != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Use current values as defaults, override with input if provided
		longURL := currentURL.LongUrl
		if input.LongURL != nil {
			longURL = *input.LongURL
		}

		expiresAt := currentURL.ExpiresAt
		if input.ExpiresAt != nil {
			expiresAt = pgtype.Timestamp{Time: *input.ExpiresAt, Valid: true}
		}

		clickLimit := currentURL.ClickLimit
		if input.ClickLimit != nil {
			clickLimit = pgtype.Int4{Int32: int32(*input.ClickLimit), Valid: true}
		}

		updatedURL, err := db.UpdateURL(r.Context(), sqlc.UpdateURLParams{
			ShortID:    shortID,
			LongUrl:    longURL,
			ExpiresAt:  expiresAt,
			ClickLimit: clickLimit,
			UserID:     sqlc.UUIDToNullable(&userID),
		})

		if err != nil {
			http.Error(w, "Failed to update URL or URL not found", http.StatusInternalServerError)
			return
		}

		// Convert response to user-friendly format
		response := map[string]interface{}{
			"short_id":   updatedURL.ShortID,
			"long_url":   updatedURL.LongUrl,
			"created_at": updatedURL.CreatedAt.Time,
		}

		if updatedURL.ExpiresAt.Valid {
			response["expires_at"] = updatedURL.ExpiresAt.Time
		}

		if updatedURL.ClickLimit.Valid {
			response["click_limit"] = updatedURL.ClickLimit.Int32
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
