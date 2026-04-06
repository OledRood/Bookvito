package main

import (
	"bookvito/config"
	"bookvito/internal/domain"
	"bookvito/pkg/database"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if cfg.AppEnv != "test" && !strings.EqualFold(os.Getenv("ALLOW_TEST_SEED"), "true") {
		log.Fatalf("refusing to seed non-test environment: APP_ENV=%s", cfg.AppEnv)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		log.Fatalf("create extension: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := db.Exec(`TRUNCATE TABLE reports, book_movement_histories, exchanges, reviews, books, users, locations RESTART IDENTITY CASCADE`).Error; err != nil {
		log.Fatalf("truncate: %v", err)
	}

	locationID := uuid.New()
	location := &domain.Location{
		ID:      locationID,
		Name:    "Main Pickup Point",
		Address: "Test street 1",
	}
	if err := db.Create(location).Error; err != nil {
		log.Fatalf("create location: %v", err)
	}

	admin := createUser("admin@bookvito.test", "Admin", domain.RoleAdmin)
	moder := createUser("moder@bookvito.test", "Moder", domain.RoleModer)
	user := createUser("user@bookvito.test", "User", domain.RoleUser)

	for _, entity := range []*domain.User{admin, moder, user} {
		if err := db.Create(entity).Error; err != nil {
			log.Fatalf("create user %s: %v", entity.Email, err)
		}
	}

	availableBook := &domain.Book{
		ID:                uuid.New(),
		OwnerID:           admin.ID,
		Title:             "Seeded Available Book",
		Author:            "Seed Author",
		Description:       "Available for request",
		Condition:         domain.ConditionGood,
		Status:            domain.BookAvailable,
		CurrentLocationID: &locationID,
	}
	reportedBook := &domain.Book{
		ID:                uuid.New(),
		OwnerID:           user.ID,
		Title:             "Reported Book",
		Author:            "Reported Author",
		Description:       "Has pending report",
		Condition:         domain.ConditionGood,
		Status:            domain.BookAvailable,
		CurrentLocationID: &locationID,
	}
	for _, book := range []*domain.Book{availableBook, reportedBook} {
		if err := db.Create(book).Error; err != nil {
			log.Fatalf("create book %s: %v", book.Title, err)
		}
	}

	report := &domain.Report{
		ID:     uuid.New(),
		BookID: reportedBook.ID,
		UserID: admin.ID,
		Reason: "Seeded report for moderation flow",
		Status: domain.ReportPending,
	}
	if err := db.Create(report).Error; err != nil {
		log.Fatalf("create report: %v", err)
	}

	if err := writeSeedOutput(cfg, admin, moder, user); err != nil {
		log.Fatalf("write seed output: %v", err)
	}
}

func createUser(email, name string, role domain.UserRole) *domain.User {
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultSeedPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	return &domain.User{
		ID:                    uuid.New(),
		Email:                 email,
		Password:              string(hash),
		Name:                  name,
		Role:                  role,
		Avatar:                "avatar1.png",
		RefreshToken:          uuid.NewString(),
		RefreshTokenExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
}

func defaultSeedPassword() string {
	value := strings.TrimSpace(os.Getenv("E2E_SEED_PASSWORD"))
	if value == "" {
		return "password123"
	}
	return value
}

type seededAuthState struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

func writeSeedOutput(cfg *config.Config, users ...*domain.User) error {
	outputPath := strings.TrimSpace(os.Getenv("E2E_SEED_OUTPUT"))
	if outputPath == "" {
		return nil
	}

	payload := map[string]seededAuthState{}
	for _, user := range users {
		accessToken, _, err := tokenPair(cfg.JWTSecret, user)
		if err != nil {
			return err
		}
		payload[string(user.Role)] = seededAuthState{
			AccessToken:  accessToken,
			RefreshToken: user.RefreshToken,
			User: &domain.User{
				ID:     user.ID,
				Email:  user.Email,
				Name:   user.Name,
				Avatar: user.Avatar,
				Role:   user.Role,
			},
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}

	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, body, 0o644)
}

func tokenPair(secret string, user *domain.User) (string, string, error) {
	claims := jwt.MapClaims{
		"sub":    user.ID,
		"userId": user.ID.String(),
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role,
		"exp":    time.Now().Add(10 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, user.RefreshToken, nil
}
