//go:build integration

package integration_test

import (
	"bookvito/config"
	"bookvito/internal/app"
	"bookvito/internal/domain"
	"bookvito/internal/repository/postgres"
	"bookvito/pkg/database"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func testConfig() *config.Config {
	return &config.Config{
		DBHost:               envOr("TEST_DB_HOST", "localhost"),
		DBPort:               envOr("TEST_DB_PORT", "15432"),
		DBUser:               envOr("TEST_DB_USER", "postgres"),
		DBPassword:           envOr("TEST_DB_PASSWORD", "postgres"),
		DBName:               envOr("TEST_DB_NAME", "bookvito_test"),
		DBSSLMode:            envOr("TEST_DB_SSLMODE", "disable"),
		ServerPort:           envOr("TEST_SERVER_PORT", "18080"),
		AppEnv:               "test",
		JWTSecret:            envOr("TEST_JWT_SECRET", "test-secret"),
		SiteURL:              envOr("TEST_SITE_URL", "http://localhost:3000"),
		GoogleBooksAPIKey:    envOr("TEST_GOOGLE_BOOKS_API_KEY", "test-key"),
		GoogleBooksBaseURL:   envOr("TEST_GOOGLE_BOOKS_BASE_URL", "http://localhost:18081"),
		StorageDriver:        envOr("TEST_STORAGE_DRIVER", "s3"),
		StorageBucket:        envOr("TEST_STORAGE_BUCKET", "book-images-test"),
		StorageEndpoint:      envOr("TEST_STORAGE_ENDPOINT", "localhost:19000"),
		StoragePublicBaseURL: envOr("TEST_STORAGE_PUBLIC_BASE_URL", "http://localhost:19000/book-images-test/"),
		StorageAccessKey:     envOr("TEST_STORAGE_ACCESS_KEY", "minioadmin"),
		StorageSecretKey:     envOr("TEST_STORAGE_SECRET_KEY", "minioadmin"),
		StorageUseSSL:        false,
		StorageLocalDir:      envOr("TEST_STORAGE_LOCAL_DIR", "./data/images-test"),
	}
}

func setupIntegrationTest(t *testing.T, cfg *config.Config) (*gorm.DB, *app.Dependencies, *gin.Engine) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		t.Fatalf("create extension: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		t.Fatalf("migrate db: %v", err)
	}

	deps, err := app.BuildDependencies(cfg, db)
	if err != nil {
		t.Fatalf("build deps: %v", err)
	}

	cleanupState(t, db, cfg)
	t.Cleanup(func() {
		cleanupState(t, db, cfg)
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db, deps, app.BuildRouter(cfg, deps)
}

func cleanupState(t *testing.T, db *gorm.DB, cfg *config.Config) {
	t.Helper()

	if err := db.Exec(`TRUNCATE TABLE reports, book_movement_histories, exchanges, reviews, books, users, locations RESTART IDENTITY CASCADE`).Error; err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	if strings.EqualFold(cfg.StorageDriver, "s3") || strings.EqualFold(cfg.StorageDriver, "minio") {
		ctx := context.Background()
		client, err := minio.New(cfg.StorageEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
			Secure: cfg.StorageUseSSL,
		})
		if err != nil {
			t.Fatalf("create minio client: %v", err)
		}
		for object := range client.ListObjects(ctx, cfg.StorageBucket, minio.ListObjectsOptions{Recursive: true}) {
			if object.Err != nil {
				t.Fatalf("list object: %v", object.Err)
			}
			if err := client.RemoveObject(ctx, cfg.StorageBucket, object.Key, minio.RemoveObjectOptions{}); err != nil {
				t.Fatalf("remove object: %v", err)
			}
		}
	}
}

func createUserWithRole(t *testing.T, db *gorm.DB, email, password, role string) *domain.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hash),
		Name:     strings.Split(email, "@")[0],
		Role:     domain.UserRole(role),
		Avatar:   "avatar1.png",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func authToken(t *testing.T, router http.Handler, email, password string) string {
	t.Helper()
	body := map[string]string{"email": email, "password": password}
	resp := performJSON(t, router, http.MethodPost, "/api/v1/users/login", body, "")
	assertStatus(t, resp, http.StatusOK)

	var payload struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if payload.AccessToken == "" {
		t.Fatalf("expected access token in login response")
	}
	return payload.AccessToken
}

func performJSON(t *testing.T, router http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func performMultipartImageUpload(t *testing.T, router http.Handler, path, token string) *httptest.ResponseRecorder {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", `form-data; name="image"; filename="cover.png"`)
	header.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatalf("create multipart part: %v", err)
	}
	if _, err := part.Write([]byte("png-image-content")); err != nil {
		t.Fatalf("write multipart body: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, path, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func decodeJSON[T any](t *testing.T, body io.Reader) T {
	t.Helper()
	var payload T
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	return payload
}

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if recorder.Code != expected {
		t.Fatalf("expected status %d, got %d, body=%s", expected, recorder.Code, recorder.Body.String())
	}
}

func envOr(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func createBookImageFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cover.png")
	if err := os.WriteFile(path, []byte("fixture"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func userByEmail(t *testing.T, db *gorm.DB, email string) *domain.User {
	t.Helper()
	repo := postgres.NewUserRepository(db)
	user, err := repo.GetByEmail(email)
	if err != nil {
		t.Fatalf("get user by email: %v", err)
	}
	return user
}
