package middleware

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sayanmondal31/api-gateway/cache"
)

// go:embed directive to load the lua file content at compile-time

//go:embed rate_limit.lua
var rateLimitLua string

var luaSHA string

// local token bucket represents the in-memory cache on this gateway pod
type localTokenBucket struct {
	mu        sync.Mutex
	tokens    int64
	expiresAt time.Time
}

var localTokenBucketMu sync.RWMutex
var localBuckets = make(map[string]*localTokenBucket)

// clean up expired buckets in background every 5 min interval
func init() { // this will start at first
	go func() {
		for {

			time.Sleep(5 * time.Minute)
			// lock for preventing updates
			localTokenBucketMu.Lock()
			for ip, bucket := range localBuckets {
				bucket.mu.Lock()
				// delete ip which is expired or token which empty or 0
				if time.Now().After(bucket.expiresAt) || bucket.tokens <= 0 {
					delete(localBuckets, ip)
				}
				bucket.mu.Unlock()

			}
			localTokenBucketMu.Unlock()
		}

	}()
}

func getLocalBucket(clientIP string) *localTokenBucket {
	localTokenBucketMu.RLock()
	bucket, exists := localBuckets[clientIP]
	localTokenBucketMu.RUnlock()

	if exists {
		return bucket
	}

	// if not exists we have to create
	localTokenBucketMu.Lock()
	// double check to prevent race condition during parallel creation
	bucket, exists = localBuckets[clientIP]

	if !exists {
		// if local cache exceed 100_000 entries, clear , to prevent mem leaks
		if len(localBuckets) > 100000 {
			// resetting...
			localBuckets = make(map[string]*localTokenBucket)
		}
		// create empty bucket
		bucket = &localTokenBucket{}
		localBuckets[clientIP] = bucket
	}
	localTokenBucketMu.Unlock()

	return bucket

}

func LoadRateLimitScript() error {
	sha, err := cache.RedisClient.ScriptLoad(context.Background(), rateLimitLua).Result()

	if err != nil {
		return fmt.Errorf("Failed to load lua script to redis: %w", err)
	}

	luaSHA = sha
	fmt.Printf("Rate limit Lua script loaded into Redis. SHA: %s\n", luaSHA)
	return nil
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

			// get the bucket based in client ip
			bucket := getLocalBucket(clientIP)

			// thread safe lock for specific client's local bucket
			bucket.mu.Lock()
			defer bucket.mu.Unlock()

			// fetch local token
			if bucket.tokens > 0 && time.Now().Before(bucket.expiresAt) {
				bucket.tokens--
				w.Header().Set("X-RateLimit-Duration-Us", "0")
				next.ServeHTTP(w, r)
				return
			}

			batchSize := 20 // 20 tokens at once
			now := time.Now().Unix()
			keyarr := []string{key}
			args := []any{capacity, refillRate, now, batchSize}

			// --- start timer
			startTime := time.Now()

			// Attempt using SHA-1 hash (fast path)
			res, err := cache.RedisClient.EvalSha(ctx, luaSHA, keyarr, args...).Result()
			if err != nil {
				// Fallback to sending full script if script is missing in Redis cache
				if strings.Contains(err.Error(), "NOSCRIPT") {

					// if redis restarts then there will be no lua script, that' why load script
					_ = LoadRateLimitScript()
					// execute
					res, err = cache.RedisClient.Eval(ctx, rateLimitLua, keyarr, args...).Result()

				}
			}

			// --- 2. MEASURE DURATION ---
			duration := time.Since(startTime)
			w.Header().Set("X-RateLimit-Duration-Us", fmt.Sprintf("%d", duration.Microseconds()))

			// fallback: if redis fails
			if err != nil {
				fmt.Printf("Redis rate limiter error: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}
			allocated, ok := res.(int64)
			if !ok || allocated <= 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Too many requests. Please try again later."}`))
				return
			}

			// store tokens from redis
			// deduct 1 token for current req and save in memory for 1 second
			bucket.tokens = allocated - 1
			bucket.expiresAt = time.Now().Add(1 * time.Second)

			next.ServeHTTP(w, r)

		})
	}
}
