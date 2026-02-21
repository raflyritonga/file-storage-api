package config

import (
	"os"
	"strconv"

    "github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr     string
	RedisPassword string
	RedisDB int

	JWTAccessSecret  string
    JWTRefreshSecret string
    JWTIssuer        string
    JWTAccessTTLMin  int
    JWTRefreshTTLHr  int
}

func Load() Config {
    loadEnvFile()

	return Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "file_storage_api"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

        RedisAddr:     getEnv("REDIS_ADDR", ""),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       getEnvInt("REDIS_DB", 0),

        JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
        JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
        JWTIssuer:        getEnv("JWT_ISSUER", ""),
        JWTAccessTTLMin:  getEnvInt("JWT_ACCESS_TTL_MIN", 15),
        JWTRefreshTTLHr:  getEnvInt("JWT_REFRESH_TTL_HOUR", 168),
	}
}

func loadEnvFile() {
    paths := []string{".env", "../.env"}
    for _, path := range paths {
        if err := godotenv.Load(path); err == nil {
            return
        }
    }
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
    value := os.Getenv(key)
    if value == "" {
        return fallback
    }

    number, err := strconv.Atoi(value)
    if err != nil {
        return fallback
    }
    return number
}