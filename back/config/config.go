package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	ServerPort        string
	JWTSecret         string
	SiteURL           string
	GoogleBooksAPIKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	loadDotEnv()

	cfg := &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "bookvito"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-here"),
		SiteURL:           getEnv("SITE_URL", "https://bookvito.ru"),
		GoogleBooksAPIKey: getEnv("GOOGLE_BOOKS_API_KEY", ""),
	}

	return cfg, nil
}

func loadDotEnv() {
	loadIfExists(".env")

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	projectEnvPath := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".env"))
	if projectEnvPath == ".env" {
		return
	}

	loadIfExists(projectEnvPath)
}

func loadIfExists(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}

	_ = godotenv.Load(path)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
