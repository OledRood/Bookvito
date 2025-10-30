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
			// Update current authenticated user (e.g., avatar/name)
			authed.PUT("/me", userHandler.UpdateMe)
			// Delete current authenticated user
			authed.DELETE("/me", userHandler.DeleteMe)
			// TODO: получить все брони, историю обменов и т.д.

		}

		books := api.Group("/books")
		{
			bookHandler := NewBookHandler(bookUC)

			// Summary is a public endpoint but we want to apply optional auth so that
			// when the client sends a valid token we can exclude user's own/returned books.
			books.GET("/summary", OptionalAuthMiddleware(cfg.JWTSecret), bookHandler.GetSummaryList)
			// Search is public but applies same visibility rules when optional auth is present
			books.GET("/search", OptionalAuthMiddleware(cfg.JWTSecret), bookHandler.Search)
			books.GET("/list", bookHandler.GetList)
			books.GET("/:id", bookHandler.GetByID)

			// Защищенные маршруты (требуют токен)
			authed := books.Group("/")
			authed.Use(AuthMiddleware(cfg.JWTSecret))
			// Upload image for books (multipart/form-data, field name: "image")
			authed.POST("/upload", bookHandler.UploadImage)
			// Attach existing uploaded image to a book: PUT /api/v1/books/image/:id
			authed.PUT("/image/:id", bookHandler.SetImage)
			authed.POST("/create", bookHandler.Create)
			authed.POST("/request", bookHandler.Request)
			// User-specific lists
			authed.GET("/my", bookHandler.GetMyBooks)
			authed.GET("/my/stats", bookHandler.GetMyBooksStats)
			authed.GET("/my/stats/:bookId", bookHandler.GetBookStats)
			authed.GET("/reserved", bookHandler.GetReservedBooks)
			authed.GET("/shelf", bookHandler.GetShelfBooks)
			authed.GET("/read", bookHandler.GetReadBooks)
			// Reservation management endpoints used by frontend for reserved books page
			authed.PUT("/reservation/extend", bookHandler.ExtendReservation)
			authed.DELETE("/reservation/cancel", bookHandler.CancelReservation)
			authed.PUT("/borrow", bookHandler.Borrow)
			authed.PUT("/return", bookHandler.Return)
			authed.DELETE("/delete", bookHandler.Delete)
		}
		locations := api.Group("/locations")
		{
			locationHandler := NewLocationHandler(locationUC)
			locations.GET("/:id", locationHandler.GetByID)
			locations.POST("/create", locationHandler.Create)
			locations.GET("/getAll", locationHandler.GetAll)

		}
	}
}
