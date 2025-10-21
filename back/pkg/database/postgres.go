package database

import (
	"bookvito/config"
	"bookvito/internal/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Создает новое подключение к базе данных PostgreSQL
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Автоматическая миграция схемы базы данных
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Location{},
		&domain.User{},
		&domain.Book{},
		&domain.Exchange{},
		&domain.Review{},
		&domain.BookMovementHistory{},
	)
}
