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
			// TODO: изменение пароля

			authed := users.Group("/")
			authed.Use(AuthMiddleware(cfg.JWTSecret))
			authed.GET("/me", userHandler.GetByID)
			// TODO: получить все брони, историю обменов и т.д.

		}

		books := api.Group("/books")
		{
			bookHandler := NewBookHandler(bookUC)

			books.GET("/summary", bookHandler.GetSummaryList)
			books.GET("/list", bookHandler.GetList)
			books.GET("/:id", bookHandler.GetByID)

			// Защищенные маршруты (требуют токен)
			authed := books.Group("/")
			authed.Use(AuthMiddleware(cfg.JWTSecret))
			authed.POST("/create", bookHandler.Create)
			authed.POST("/request", bookHandler.Request)
			authed.PUT("/borrow", bookHandler.Borrow)
			authed.PUT("/return", bookHandler.Return)
			authed.DELETE("/delete", bookHandler.Delete)
		}
		locations := api.Group("/locations")
		{
			locationHandler := NewLocationHandler(locationUC)
			locations.GET("/:id", locationHandler.GetByID)
			locations.GET("/getAll", locationHandler.GetAll)

		}
	}
}
