package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayanmondal31/api-gateway/cache"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type UserClaims struct {
	UserId string `json:"sub"`
	Email  string `json:"email"`
}

// Authenticate verifies JWTs, check redis for blacklisted users, inject headers/context
func Authenticate(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// clear spoof headers send by untrusted client
			r.Header.Del("X-USER-ID")
			r.Header.Del("X-USER-EMAIL")

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, `{"error":"Missing Authorization Header"}`, http.StatusUnauthorized)
				return
			}

			partsToken := strings.Split(authHeader, " ")

			if partsToken[0] != "Bearer" {
				http.Error(w, `{"error":"Invalid Authentication format"}`, http.StatusUnauthorized)
				return
			}

			tokenString := partsToken[1]

			// check redis blacklist
			isBlackListed, err := cache.RedisClient.Exists(r.Context(), "blacklist:"+tokenString).Result()

			if err == nil && isBlackListed > 0 {
				http.Error(w, `{"error":"Token has been revoked!"}`, http.StatusUnauthorized)
				return
			}

			// 3. Cryptographic Signature Validation
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"Invalid token claims"}`, http.StatusUnauthorized)
				return
			}
			userID, _ := claims["sub"].(string)
			email, _ := claims["email"].(string)
			if userID == "" {
				http.Error(w, `{"error":"Invalid subject claim"}`, http.StatusUnauthorized)
				return
			}
			// 4. Inject claims into HTTP headers for downstream microservices
			r.Header.Set("X-User-Id", userID)
			r.Header.Set("X-User-Email", email)
			// 5. Inject claims into Go Request Context for internal Gateway use
			ctx := context.WithValue(r.Context(), UserContextKey, UserClaims{
				UserId: userID,
				Email:  email,
			})
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
