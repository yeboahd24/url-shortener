package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			key := "rate:" + ip
			count, _ := redisClient.Get(r.Context(), key).Int()
			if count >= 10 {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			redisClient.Incr(r.Context(), key)
			redisClient.Expire(r.Context(), key, time.Minute)
			next.ServeHTTP(w, r)
		})
	}
}
