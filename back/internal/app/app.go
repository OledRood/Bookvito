package app

import (
	"bookvito/config"
	delivery "bookvito/internal/delivery/http"
	"bookvito/internal/domain"
	"bookvito/internal/repository/postgres"
	"bookvito/internal/service/googlebooks"
	"bookvito/internal/service/storage"
	"bookvito/internal/usecase"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Dependencies struct {
	UserUseCase     domain.UserUseCase
	BookUseCase     domain.BookUseCase
	ExchangeUseCase domain.ExchangeUseCase
	LocationUseCase domain.LocationUseCase
	AdminUseCase    domain.AdminUseCase
	ModerUseCase    domain.ModerUseCase
	ImageStorage    domain.ImageStorage
}

func BuildDependencies(cfg *config.Config, db *gorm.DB) (*Dependencies, error) {
	userRepo := postgres.NewUserRepository(db)
	bookRepo := postgres.NewBookRepository(db)
	exchangeRepo := postgres.NewExchangeRepository(db)
	movementRepo := postgres.NewBookMovementHistoryRepository(db)
	locationRepo := postgres.NewLocationRepository(db)
	reportRepo := postgres.NewReportRepository(db)

	imageStorage, err := storage.NewImageStorage(
		cfg.StorageDriver,
		cfg.StorageLocalDir,
		cfg.StoragePublicBaseURL,
		cfg.StorageEndpoint,
		cfg.StorageAccessKey,
		cfg.StorageSecretKey,
		cfg.StorageBucket,
		cfg.StorageUseSSL,
	)
	if err != nil {
		return nil, err
	}

	googleBooksSvc := googlebooks.New(cfg.GoogleBooksAPIKey, cfg.GoogleBooksBaseURL, nil)

	return &Dependencies{
		UserUseCase:     usecase.NewUserUseCase(userRepo, movementRepo, cfg.JWTSecret),
		BookUseCase:     usecase.NewBookUseCase(bookRepo, movementRepo, exchangeRepo, locationRepo, googleBooksSvc, imageStorage),
		ExchangeUseCase: usecase.NewExchangeUseCase(exchangeRepo, bookRepo, userRepo, movementRepo),
		LocationUseCase: usecase.NewLocationUseCase(locationRepo),
		AdminUseCase:    usecase.NewAdminUseCase(userRepo, bookRepo, exchangeRepo, reportRepo),
		ModerUseCase:    usecase.NewModerUseCase(reportRepo, bookRepo),
		ImageStorage:    imageStorage,
	}, nil
}

func BuildRouter(cfg *config.Config, deps *Dependencies) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}

			normalizedOrigin := strings.TrimRight(strings.ToLower(origin), "/")
			normalizedSiteURL := strings.TrimRight(strings.ToLower(cfg.SiteURL), "/")

			return normalizedOrigin == normalizedSiteURL || normalizedOrigin == "https://www.bookvito.ru"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if strings.EqualFold(cfg.StorageDriver, "local") || cfg.StorageDriver == "" {
		router.Static("/images", cfg.StorageLocalDir)
	}

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			delivery.WriteError(c, 404, domain.ErrorCodeNotFound, "Маршрут не найден")
			return
		}
		delivery.WriteError(c, 404, domain.ErrorCodeNotFound, "Страница не найдена")
	})

	router.NoMethod(func(c *gin.Context) {
		delivery.WriteError(c, 405, domain.ErrorCodeValidation, "Метод не поддерживается")
	})

	delivery.NewRouter(
		router,
		deps.UserUseCase,
		deps.BookUseCase,
		deps.ExchangeUseCase,
		deps.LocationUseCase,
		deps.AdminUseCase,
		deps.ModerUseCase,
		deps.ImageStorage,
		cfg,
	)

	return router
}
