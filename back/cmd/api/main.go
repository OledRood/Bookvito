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
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	bookRepo := postgres.NewBookRepository(db)
	exchangeRepo := postgres.NewExchangeRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	bookUseCase := usecase.NewBookUseCase(bookRepo)
	exchangeUseCase := usecase.NewExchangeUseCase(exchangeRepo, bookRepo, userRepo)

	// Initialize HTTP handlers
	router := gin.Default()
	http.NewRouter(router, userUseCase, bookUseCase, exchangeUseCase)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
