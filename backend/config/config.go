package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port            string
	DBConnString    string
	DBName          string
	JWTSecret       string
	MigrationsPath  string
	AllowedOrigins  []string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

func LoadConfig() *Config {
	dbUser := getEnv("DB_USER")
	dbPass := getEnv("DB_PASSWORD")
	dbName := getEnv("DB_NAME")
	jwtSecret := getEnv("JWT_SECRET")
	migrationsPath := getEnvWithFallback("MIGRATIONS_PATH", "migrations")
	port := getEnvWithFallback("PORT", "8080")
	dbHost := getEnvWithFallback("DB_HOST", "localhost")
	dbPort := getEnvWithFallback("DB_PORT", "5433")
	allowedOrigins := getAllowedOrigins()
	accessTokenDuration := GetDurationEnv("ACCESS_TOKEN_TTL_MINUTES", 15, time.Minute)
	refreshTokenDuration := GetDurationEnv("REFRESH_TOKEN_TTL_HOURS", 24*7, time.Hour)

	if dbHost == "db" {
		if _, err := os.Stat("/.dockerenv"); err != nil {
			dbHost = "localhost"
		}
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	return &Config{
		Port:            port,
		DBConnString:    dbUrl,
		DBName:          dbName,
		JWTSecret:       jwtSecret,
		MigrationsPath:  migrationsPath,
		AllowedOrigins:  allowedOrigins,
		AccessTokenTTL:  accessTokenDuration,
		RefreshTokenTTL: refreshTokenDuration,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     60 * time.Second,
	}
}

func getAllowedOrigins() []string {
	raw := getEnvWithFallback("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")
	parts := strings.Split(raw, ",")
	allowedOrigins := make([]string, 0, len(parts))

	for _, origin := range parts {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins = append(allowedOrigins, origin)
		}
	}

	if len(allowedOrigins) == 0 {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	return allowedOrigins
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("FATAL: Environment variable %s is missing", key)
	}
	return value
}

func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetDurationEnv(key string, fallback int, unit time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallback) * unit
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(fallback) * unit
	}
	return time.Duration(parsed) * unit
}
