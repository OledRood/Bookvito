package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	ServerPort           string
	AppEnv               string
	JWTSecret            string
	SiteURL              string
	GoogleBooksAPIKey    string
	GoogleBooksBaseURL   string
	StorageDriver        string
	StorageBucket        string
	StorageEndpoint      string
	StoragePublicBaseURL string
	StorageAccessKey     string
	StorageSecretKey     string
	StorageUseSSL        bool
	StorageLocalDir      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	loadDotEnv()

	cfg := &Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "bookvito"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		AppEnv:               getEnv("APP_ENV", "dev"),
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key-here"),
		SiteURL:              getEnv("SITE_URL", "https://bookvito.ru"),
		GoogleBooksAPIKey:    getEnv("GOOGLE_BOOKS_API_KEY", ""),
		GoogleBooksBaseURL:   getEnv("GOOGLE_BOOKS_BASE_URL", "https://www.googleapis.com/books/v1"),
		StorageDriver:        getEnv("STORAGE_DRIVER", "local"),
		StorageBucket:        getEnv("STORAGE_BUCKET", "book-images"),
		StorageEndpoint:      getEnv("STORAGE_ENDPOINT", "localhost:9000"),
		StoragePublicBaseURL: getEnv("STORAGE_PUBLIC_BASE_URL", "/images/"),
		StorageAccessKey:     getEnv("STORAGE_ACCESS_KEY", "minioadmin"),
		StorageSecretKey:     getEnv("STORAGE_SECRET_KEY", "minioadmin"),
		StorageUseSSL:        getEnvAsBool("STORAGE_USE_SSL", false),
		StorageLocalDir:      getEnv("STORAGE_LOCAL_DIR", "./data/images"),
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

func getEnvAsBool(key string, defaultValue bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}

	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}
