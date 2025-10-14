package http
package http

import (
	"bookvito/internal/usecase"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes all routes
func NewRouter(router *gin.Engine, userUC *usecase.UserUseCase, bookUC *usecase.BookUseCase, exchangeUC *usecase.ExchangeUseCase) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		{
			userHandler := NewUserHandler(userUC)
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.GET("", userHandler.List)
		}

		// Book routes
		books := api.Group("/books")
		{
			bookHandler := NewBookHandler(bookUC)
			books.POST("", bookHandler.Create)
			books.GET("/:id", bookHandler.GetByID)
			books.GET("", bookHandler.List)
			books.GET("/search", bookHandler.Search)
			books.GET("/available", bookHandler.GetAvailable)
			books.PUT("/:id", bookHandler.Update)
			books.DELETE("/:id", bookHandler.Delete)
			books.GET("/owner/:owner_id", bookHandler.GetByOwner)
		}

		// Exchange routes
		exchanges := api.Group("/exchanges")
		{
			exchangeHandler := NewExchangeHandler(exchangeUC)
			exchanges.POST("", exchangeHandler.Create)
			exchanges.GET("/:id", exchangeHandler.GetByID)
			exchanges.GET("", exchangeHandler.List)
			exchanges.GET("/requester/:requester_id", exchangeHandler.GetByRequester)
			exchanges.GET("/owner/:owner_id", exchangeHandler.GetByOwner)
			exchanges.PUT("/:id/accept", exchangeHandler.Accept)
			exchanges.PUT("/:id/reject", exchangeHandler.Reject)
			exchanges.PUT("/:id/complete", exchangeHandler.Complete)
		}
	}
}
