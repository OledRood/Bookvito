package http

import (
	"bookvito/internal/domain"
	"net/http"

	// "strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookHandler struct {
	bookUC domain.BookUseCase
}

func NewBookHandler(bookUC domain.BookUseCase) *BookHandler {
	return &BookHandler{bookUC: bookUC}
}

type CreateBookRequest struct {
	Title             string               `json:"title" binding:"required"`
	Author            string               `json:"author" binding:"required"`
	Description       string               `json:"description"`
	Condition         domain.BookCondition `json:"condition" binding:"required,oneof=excellent good bad"`
	ImageURL          string               `json:"image_url"`
	CurrentLocationID *uuid.UUID           `json:"current_location_id"`
}

func (h *BookHandler) Create(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, exists := c.Get("userId")
	println("UserID from context:", userIdRaw)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
		return
	}

	err := h.bookUC.CreateBook(&domain.Book{
		Title:             req.Title,
		Author:            req.Author,
		Description:       req.Description,
		Condition:         req.Condition,
		ImageURL:          req.ImageURL,
		CurrentLocationID: req.CurrentLocationID,
		OwnerID:           uuid.MustParse(userIDStr),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book created successfully"})
}

func (h *BookHandler) GetSummaryList(c *gin.Context) {
	books, err := h.bookUC.GetSummaryBooksList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetList(c *gin.Context) {
	books, err := h.bookUC.GetBooksList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	book, err := h.bookUC.GetBookByID(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	//Здесь нужно получить userID из токена
	//Проверить его валидность и идти дальше если все хорошо

	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
		return
	}

	err = h.bookUC.DeleteBook(bookID, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}
func (h *BookHandler) GetBookMovementHistory(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	history, err := h.bookUC.GetBookMovementHistory(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
