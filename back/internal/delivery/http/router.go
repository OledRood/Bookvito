package http

import (
	"bookvito/config"
	"bookvito/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	router *gin.Engine,
	userUC domain.UserUseCase,
	bookUC domain.BookUseCase,
	exchangeUC domain.ExchangeUseCase,
	locationUC domain.LocationUseCase,
	adminUC domain.AdminUseCase,
	moderUC domain.ModerUseCase,
	cfg *config.Config,
) {
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
			// Logout: revokes refresh token in DB
			authed.POST("/logout", userHandler.Logout)
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
			locations.GET("/getAll", locationHandler.GetAll)

			// Защищённые маршруты локаций (только admin)
			locationsAuthed := locations.Group("/")
			locationsAuthed.Use(AuthMiddleware(cfg.JWTSecret), RequireRole("admin"))
			locationsAuthed.POST("/create", locationHandler.Create)
			locationsAuthed.PUT("/:id", locationHandler.Update)
			locationsAuthed.DELETE("/:id", locationHandler.Delete)
		}

		// Маршруты книги — жалоба от авторизованного пользователя
		moderHandler := NewModerHandler(moderUC)
		api.POST("/books/:id/report", AuthMiddleware(cfg.JWTSecret), moderHandler.CreateReport)

		// Маршруты модератора (moder + admin)
		moder := api.Group("/moder")
		moder.Use(AuthMiddleware(cfg.JWTSecret), RequireRole("moder", "admin"))
		{
			moder.GET("/reports", moderHandler.GetReports)
			moder.PUT("/reports/:id/resolve", moderHandler.ResolveReport)
			moder.PUT("/reports/:id/dismiss", moderHandler.DismissReport)
			moder.PUT("/books/:id/archive", moderHandler.ArchiveBook)
		}

		// Маршруты администратора (только admin)
		adminHandler := NewAdminHandler(userUC, adminUC)
		admin := api.Group("/admin")
		admin.Use(AuthMiddleware(cfg.JWTSecret), RequireRole("admin"))
		{
			admin.GET("/stats", adminHandler.GetStats)
			admin.GET("/users", adminHandler.ListUsers)
			admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		}
	}
}
