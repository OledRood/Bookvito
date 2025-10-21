package http

import (
	"bookvito/config"
	"bookvito/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, userUC domain.UserUseCase, bookUC domain.BookUseCase, exchangeUC domain.ExchangeUseCase, locationUC domain.LocationUseCase, cfg *config.Config) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{

		users := api.Group("/users")
		{
			userHandler := NewUserHandler(userUC)
			users.POST("/registration", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/refresh", userHandler.Refresh)
			// 		users.GET("/:id", userHandler.GetByID)
			// 		users.PUT("/:id", userHandler.Update)
			// 		users.DELETE("/:id", userHandler.Delete)
			// 		users.GET("", userHandler.List)
		}

		// --- Маршруты для книг ---
		books := api.Group("/books")
		{
			bookHandler := NewBookHandler(bookUC)
			// Публичные маршруты
			books.GET("/summary", bookHandler.GetSummaryList)
			books.GET("/list", bookHandler.GetList)
			books.GET("/:id", bookHandler.GetByID)

			// Защищенные маршруты (требуют токен)
			authed := books.Group("/")
			authed.Use(AuthMiddleware(cfg.JWTSecret))
			authed.POST("/creare", bookHandler.Create)
			// authed.POST("/:id/reserve", bookHandler.Reserve) // TODO: Implement Reserve method in BookHandler
			// authed.DELETE("/:id", bookHandler.Delete)
		}
		userHandler := NewUserHandler(userUC)
		locationHandler := NewLocationHandler(locationUC)

		// --- Маршруты для локаций ---
		locations := api.Group("/locations")
		{
			locations.GET("/:id", locationHandler.GetByID)
			locations.GET("/getAll", locationHandler.GetAll)

			authed := users.Group("/")
			authed.Use(AuthMiddleware(cfg.JWTSecret))
			authed.GET("me", userHandler.GetByID)
			authed.GET("me/history", userHandler.GetMyMovementHistory)
			// authed.PUT("me", userHandler.Update)
			// authed.DELETE("me", userHandler.Delete)

		}
		// Защищенные маршруты (требуют токен и прав администратора)
		admin := locations.Group("/")
		admin.Use(AuthMiddleware(cfg.JWTSecret))
		admin.POST("/create", locationHandler.Create)
		admin.PUT("/:id", locationHandler.Update)
		admin.DELETE("/:id", locationHandler.Delete)
	}
}
