package config

import "os"

type Config struct {
	Port       string
	RedisURL   string
	AuthSvcURL string
	// DispatchURL string
	JWTSecret string
}

func Load() *Config {
	Port := getEnv("PORT", "3000")
	RedisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	AuthSvcURL := getEnv("AUTH_SERVICE_URL", "http://localhost:3001")
	JWTSecret := getEnv("JWT_SECRET", "")

	return &Config{
		Port:       Port,
		RedisURL:   RedisURL,
		AuthSvcURL: AuthSvcURL,
		JWTSecret:  JWTSecret,
	}

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}

	return fallback
}
