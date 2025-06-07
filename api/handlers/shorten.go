package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// ShortenURLRequest represents the request body for shortening a URL
type ShortenURLRequest struct {
	LongURL    string     `json:"long_url" example:"https://example.com" binding:"required"`
	CustomID   string     `json:"custom_id,omitempty" example:"my-custom-url"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	ClickLimit *int       `json:"click_limit,omitempty" example:"100"`
}

// ShortenURLResponse represents the response for shortening a URL
type ShortenURLResponse struct {
	ShortURL string `json:"short_url" example:"abc123"`
}

func timeToNullable(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}

func intToNullable(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*i), Valid: true}
}

// ShortenURL creates a shortened URL
// @Summary Shorten URL
// @Description Create a shortened URL. Custom IDs require authentication.
// @Tags urls
// @Accept json
// @Produce json
// @Param url body ShortenURLRequest true "URL to shorten"
// @Success 200 {object} ShortenURLResponse "URL shortened successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Authentication required for custom URLs"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /shorten [post]
// @Router /api/shorten [post]
func ShortenURL(db *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			LongURL    string     `json:"long_url"`
			CustomID   string     `json:"custom_id"`
			ExpiresAt  *time.Time `json:"expires_at"`
			ClickLimit *int       `json:"click_limit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var userID *uuid.UUID
		if input.CustomID != "" {
			uidStr, ok := r.Context().Value("user_id").(string)
			if !ok {
				http.Error(w, "Authentication required for custom URLs", http.StatusUnauthorized)
				return
			}
			uid, err := uuid.Parse(uidStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			userID = &uid
		}

		shortID := input.CustomID
		if shortID == "" {
			shortID = generateShortID()
		}

		_, err := db.CreateURL(r.Context(), sqlc.CreateURLParams{
			ShortID:    shortID,
			LongUrl:    input.LongURL,
			UserID:     sqlc.UUIDToNullable(userID),
			CreatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
			ExpiresAt:  timeToNullable(input.ExpiresAt),
			ClickLimit: intToNullable(input.ClickLimit),
		})
		if err != nil {
			http.Error(w, "Failed to create URL", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"short_url": shortID})
	}
}

func generateShortID() string {
	// Base62 characters: 0-9, a-z, A-Z
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const length = 8

	result := make([]byte, length)
	// Generate random bytes
	randomBytes := uuid.New()

	// Convert to base62
	for i := 0; i < length; i++ {
		// Use different parts of UUID for better distribution
		idx := int(randomBytes[i%16]) % len(charset)
		result[i] = charset[idx]
	}

	return string(result)
}
