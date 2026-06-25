package middleware

import (
	"crypto/sha1"
	_ "embed"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sayanmondal31/api-gateway/cache"
)

// go:embed directive to load the lua file content at compile-time

//go:embed rate_limit.lua
var rateLimitLua string

var luaSHA string

func init() {

	// Generate SHA-1 of the embedded Lua script
	h := sha1.New()
	h.Write([]byte(rateLimitLua))
	luaSHA = hex.EncodeToString(h.Sum(nil))

}

func GetClientIP(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		addresses := r.Header.Get(header)
		if addresses != "" {
			parts := strings.Split(addresses, ",")

			ip := strings.TrimSpace(parts[0])

			if ip != "" {
				return ip
			}
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}

// Read this
// RateLimit creates a middleware handler to enforce rate limits
func RateLimit(capacity int, refillRate float64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			clientIP := GetClientIP(r)
			key := fmt.Sprintf("rate_limit:%s", clientIP)
			now := time.Now().Unix()
			// Attempt using SHA-1 hash (fast path)
			res, err := cache.RedisClient.EvalSha(ctx, luaSHA, []string{key}, capacity, refillRate, now, 1).Result()
			if err != nil {
				// Fallback to sending full script if script is missing in Redis cache
				if strings.Contains(err.Error(), "NOSCRIPT") {
					res, err = cache.RedisClient.Eval(ctx, rateLimitLua, []string{key}, capacity, refillRate, now, 1).Result()
				}
			}
			// Fail-Open logic for reliability
			if err != nil {
				fmt.Printf("Redis rate limiter error: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}
			allowed, ok := res.(int64)
			if !ok || allowed != 1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Too many requests. Please try again later."}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
