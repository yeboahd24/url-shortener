package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username string `json:"username" example:"john_doe" binding:"required"`
	Email    string `json:"email" example:"john@example.com" binding:"required"`
}

// CreateUserResponse represents the response for creating a user
type CreateUserResponse struct {
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string `json:"username" example:"john_doe"`
	Email     string `json:"email" example:"john@example.com"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// CreateUser creates a new user
// @Summary Create User
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User information"
// @Success 200 {object} CreateUserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func CreateUser(db *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.Username == "" || input.Email == "" {
			http.Error(w, "Username and email are required", http.StatusBadRequest)
			return
		}

		userID := uuid.New()
		user, err := db.CreateUser(r.Context(), sqlc.CreateUserParams{
			UserID:    userID,
			Username:  input.Username,
			Email:     input.Email,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":    user.UserID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt.Time,
		})
	}
}

// CreateAPIKeyResponse represents the response for creating an API key
type CreateAPIKeyResponse struct {
	APIKey    string `json:"api_key" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// CreateAPIKey creates a new API key for the authenticated user
// @Summary Create API Key
// @Description Create a new API key for the authenticated user
// @Tags api-keys
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} CreateAPIKeyResponse "API key created successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/keys [post]
func CreateAPIKey(db *sqlc.Queries) http.HandlerFunc {
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

		apiKey := uuid.New()
		key, err := db.CreateAPIKey(r.Context(), sqlc.CreateAPIKeyParams{
			Key:       apiKey,
			UserID:    userID,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			http.Error(w, "Failed to create API key", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"api_key":    key.Key,
			"user_id":    key.UserID,
			"created_at": key.CreatedAt.Time,
		})
	}
}

// ListAPIKeysResponse represents the response for listing API keys
type ListAPIKeysResponse struct {
	APIKeys []APIKeyInfo `json:"api_keys"`
}

// APIKeyInfo represents API key information
type APIKeyInfo struct {
	Key       string `json:"key" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// ListAPIKeys lists all API keys for the authenticated user
// @Summary List API Keys
// @Description List all API keys for the authenticated user
// @Tags api-keys
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} ListAPIKeysResponse "API keys retrieved successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/keys [get]
func ListAPIKeys(db *sqlc.Queries) http.HandlerFunc {
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

		keys, err := db.ListUserAPIKeys(r.Context(), userID)
		if err != nil {
			http.Error(w, "Failed to fetch API keys", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"api_keys": keys,
		})
	}
}

// DeleteAPIKeyRequest represents the request body for deleting an API key
type DeleteAPIKeyRequest struct {
	APIKey string `json:"api_key" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
}

// DeleteAPIKey deletes an API key for the authenticated user
// @Summary Delete API Key
// @Description Delete an API key for the authenticated user
// @Tags api-keys
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param apikey body DeleteAPIKeyRequest true "API key to delete"
// @Success 200 {object} map[string]string "API key deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/keys [delete]
func DeleteAPIKey(db *sqlc.Queries) http.HandlerFunc {
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

		var input struct {
			APIKey string `json:"api_key"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		apiKey, err := uuid.Parse(input.APIKey)
		if err != nil {
			http.Error(w, "Invalid API key format", http.StatusBadRequest)
			return
		}

		err = db.DeleteAPIKey(r.Context(), sqlc.DeleteAPIKeyParams{
			Key:    apiKey,
			UserID: userID,
		})

		if err != nil {
			http.Error(w, "Failed to delete API key", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "API key deleted successfully",
		})
	}
}
