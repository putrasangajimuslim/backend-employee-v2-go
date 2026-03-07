// internal/infrastructure/config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	AppEnv   string
	DBDriver string
	DBPath   string
	DBDsn    string
	JWTSecret string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		AppPort:   getEnv("APP_PORT", "3000"),
		AppEnv:    getEnv("APP_ENV", "development"),
		DBDriver:  getEnv("DB_DRIVER", "sqlite"),
		DBPath:    getEnv("DB_PATH", "./app.db"),
		DBDsn:     getEnv("DB_DSN", ""),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
