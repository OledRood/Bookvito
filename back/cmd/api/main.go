package main

import (
	"bookvito/config"
	"bookvito/internal/delivery/http"
	"bookvito/internal/repository/postgres"
	"bookvito/internal/usecase"
	"bookvito/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Инициализируем подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Auto-migrate database schema
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Не удалось выполнить миграцию базы данных: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	bookRepo := postgres.NewBookRepository(db)
	exchangeRepo := postgres.NewExchangeRepository(db)
	movementRepo := postgres.NewBookMovementHistoryRepository(db)
	locationRepo := postgres.NewLocationRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, movementRepo, cfg.JWTSecret)
	bookUseCase := usecase.NewBookUseCase(bookRepo, movementRepo, exchangeRepo)
	exchangeUseCase := usecase.NewExchangeUseCase(exchangeRepo, bookRepo, userRepo, movementRepo)
	locationUseCase := usecase.NewLocationUseCase(locationRepo)

	// Initialize HTTP handlers
	router := gin.Default()
	http.NewRouter(router, userUseCase, bookUseCase, exchangeUseCase, locationUseCase, cfg)
	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
