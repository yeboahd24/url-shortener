package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/yeboahd24/url-shortener/api/handlers"
	"github.com/yeboahd24/url-shortener/api/middleware"
	"github.com/yeboahd24/url-shortener/config"
	_ "github.com/yeboahd24/url-shortener/docs"
	"github.com/yeboahd24/url-shortener/queries/sqlc"
)

// @title URL Shortener API
// @version 1.0
// @description A comprehensive URL shortener service with analytics, user management, and API key authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key for authentication

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgx.Connect(context.Background(), cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
	})

	queries := sqlc.New(db)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RateLimitMiddleware(redisClient))

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9000/swagger/doc.json"),
	))

	// Health and stats routes (public)
	r.Get("/health", handlers.HealthCheck(db, redisClient))
	r.Get("/stats", handlers.GetStats(queries))

	// User management routes (public)
	r.Post("/users", handlers.CreateUser(queries))

	// URL shortening (public)
	r.Post("/shorten", handlers.ShortenURL(queries))

	// Authenticated routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(queries))

		// URL shortening (authenticated - for custom URLs and advanced features)
		r.Post("/shorten", handlers.ShortenURL(queries))

		// Analytics
		r.Get("/analytics/{shortID}", handlers.GetAnalytics(queries, cfg.GeoAPIURL))

		// API key management
		r.Post("/keys", handlers.CreateAPIKey(queries))
		r.Get("/keys", handlers.ListAPIKeys(queries))
		r.Delete("/keys", handlers.DeleteAPIKey(queries))

		// URL management
		r.Get("/urls", handlers.ListUserURLs(queries))
		r.Delete("/urls/{shortID}", handlers.DeleteURL(queries))
		r.Put("/urls/{shortID}", handlers.UpdateURL(queries))
	})

	// Redirect route (must be last to avoid conflicts)
	r.Get("/{shortID}", handlers.RedirectURL(queries, redisClient, cfg.GeoAPIURL))

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
