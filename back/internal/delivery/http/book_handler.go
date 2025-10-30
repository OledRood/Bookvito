package http

import (
	"bookvito/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

// UploadImage handles multipart image upload. Field name: "image".
// It saves file into ./data/images and returns JSON { "url": "/images/<filename>" }.
func (h *BookHandler) UploadImage(c *gin.Context) {
	// single file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется файл изображения"})
		return
	}

	// validate extension
	ext := filepath.Ext(file.Filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		// ok
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неподдерживаемый формат изображения"})
		return
	}

	// ensure dir exists
	imagesDir := "./data/images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать каталог изображений"})
		return
	}

	// create unique filename
	filename := uuid.New().String() + ext
	destPath := filepath.Join(imagesDir, filename)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить файл"})
		return
	}

	// return public URL (served under /images)
	url := fmt.Sprintf("/images/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *BookHandler) GetSummaryList(c *gin.Context) {
	// If the request contains a valid authenticated user (middleware sets "userId"),
	// exclude user's own books and books they've interacted with.
	var userIDPtr *uuid.UUID
	if userIdRaw, exists := c.Get("userId"); exists {
		if userStr, ok := userIdRaw.(string); ok && userStr != "" {
			if parsed, err := uuid.Parse(userStr); err == nil {
				userIDPtr = &parsed
			}
		}
	}

	books, err := h.bookUC.GetSummaryBooksList(userIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// Search handler: GET /api/v1/books/search?q=...&limit=...&offset=...
// Uses OptionalAuthMiddleware so if a valid token is provided we can exclude
// user's own/returned books (same as summary).
func (h *BookHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if strings.TrimSpace(q) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр запроса 'q' обязателен"})
		return
	}

	limit := 100
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	var userIDPtr *uuid.UUID
	if userIdRaw, exists := c.Get("userId"); exists {
		if userStr, ok := userIdRaw.(string); ok && userStr != "" {
			if parsed, err := uuid.Parse(userStr); err == nil {
				userIDPtr = &parsed
			}
		}
	}

	books, err := h.bookUC.SearchBooks(q, limit, offset, userIDPtr)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	book, err := h.bookUC.GetBookByID(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetBookMovementHistory(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	history, err := h.bookUC.GetBookMovementHistory(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetBookStats returns statistics for a single book (requires auth)
func (h *BookHandler) GetBookStats(c *gin.Context) {
	bookIdParam := c.Param("bookId")
	if bookIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется идентификатор книги"})
		return
	}
	bookID, err := uuid.Parse(bookIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	// ensure authenticated
	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	if userIdRaw == nil {
		c.JSON(401, gin.H{"error": "Неверный пользователь"})
		return
	}

	stats, err := h.bookUC.GetBookStats(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, stats)
}

//MARK: important part

type CreateBookRequest struct {
	Title             string               `json:"title" binding:"required"`
	Author            string               `json:"author" binding:"required"`
	Description       string               `json:"description"`
	Condition         domain.BookCondition `json:"condition" binding:"required,oneof=excellent good bad"`
	ImageURL          string               `json:"image_url"`
	CurrentLocationID *uuid.UUID           `json:"current_location_id"`
}

type BookIdRequest struct {
	BookID uuid.UUID `json:"book_id" binding:"required"`
}

type DeleteBookRequest struct {
	BookID uuid.UUID `json:"book_id" binding:"required"`
	Reason string    `json:"reason"`
}

func (h *BookHandler) Create(c *gin.Context) {
	var req CreateBookRequest

	// Read raw body so we can accept both snake_case (image_url) and camelCase (imageUrl)
	// coming from different frontend builds. We reuse the body for normal binding.
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	// restore Body so binding can read it
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// try to unmarshal into a map to look for alternative keys
	var bodyMap map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &bodyMap)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fallback: if image_url was not set by binding, check camelCase keys
	if req.ImageURL == "" {
		if v, ok := bodyMap["imageUrl"].(string); ok && v != "" {
			req.ImageURL = v
		} else if v2, ok := bodyMap["image_url"].(string); ok && v2 != "" {
			req.ImageURL = v2
		}
	}

	// Fallback for current location id: accept both current_location_id (snake) and locationId (camel)
	if req.CurrentLocationID == nil {
		if v, ok := bodyMap["locationId"].(string); ok && v != "" {
			if parsed, err := uuid.Parse(v); err == nil {
				req.CurrentLocationID = &parsed
			}
		} else if v2, ok := bodyMap["current_location_id"].(string); ok && v2 != "" {
			if parsed, err := uuid.Parse(v2); err == nil {
				req.CurrentLocationID = &parsed
			}
		}
	}

	userIdRaw, exists := c.Get("userId")
	println("UserID from context:", userIdRaw.(string))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя в токене"})
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

	c.JSON(http.StatusCreated, gin.H{"message": "Книга успешно создана"})
}

func (h *BookHandler) Request(c *gin.Context) {
	var req BookIdRequest

	// Read raw body for better debug messages if binding fails
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	// restore body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&req); err != nil {
		// try tolerant parsing: accept book_id, bookId or id as string UUID
		var m map[string]interface{}
		if err2 := json.Unmarshal(bodyBytes, &m); err2 == nil {
			var idStr string
			if v, ok := m["book_id"].(string); ok && v != "" {
				idStr = v
			} else if v, ok := m["bookId"].(string); ok && v != "" {
				idStr = v
			} else if v, ok := m["id"].(string); ok && v != "" {
				idStr = v
			}
			if idStr != "" {
				if parsed, perr := uuid.Parse(idStr); perr == nil {
					req.BookID = parsed
				} else {
					fmt.Println("Request parse uuid error:", perr.Error())
					c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат идентификатора книги", "body": string(bodyBytes)})
					return
				}
			} else {
				// nothing we could extract
				fmt.Println("Request bind error:", err.Error())
				fmt.Println("Raw body:", string(bodyBytes))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": string(bodyBytes)})
				return
			}
		} else {
			// cannot even parse body as JSON
			fmt.Println("Request bind error:", err.Error())
			fmt.Println("Raw body (unparseable):", string(bodyBytes))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": string(bodyBytes)})
			return
		}
	}

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	// Tolerant parsing for optional locationId (accept both locationId and location_id)
	var locIDPtr *uuid.UUID
	var bodyMap map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &bodyMap)
	if v, ok := bodyMap["locationId"].(string); ok && v != "" {
		if parsed, perr := uuid.Parse(v); perr == nil {
			locIDPtr = &parsed
		}
	} else if v2, ok := bodyMap["location_id"].(string); ok && v2 != "" {
		if parsed, perr := uuid.Parse(v2); perr == nil {
			locIDPtr = &parsed
		}
	}

	err = h.bookUC.Request(req.BookID, userUUID, locIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно забронирована"})

}

// GetMyBooks returns books that belong to the authenticated user
func (h *BookHandler) GetMyBooks(c *gin.Context) {
	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	books, err := h.bookUC.GetBooksByOwner(userUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, books)
}

// GetMyBooksStats returns aggregated statistics for the authenticated user's books
func (h *BookHandler) GetMyBooksStats(c *gin.Context) {
	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	stats, err := h.bookUC.GetMyBooksStats(userUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, stats)
}

// GetReservedBooks returns books that the authenticated user has requested (reservations)
func (h *BookHandler) GetReservedBooks(c *gin.Context) {
	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	exchanges, err := h.bookUC.GetExchangesFromUserByStatus(userUUID, domain.ExchangeRequested)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, exchanges)
}

// ExtendReservation extends an active reservation by +1 day (subject to constraints)
func (h *BookHandler) ExtendReservation(c *gin.Context) {
	var req BookIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	if err := h.bookUC.ExtendReservation(req.BookID, userUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Бронирование продлено"})
}

// CancelReservation cancels a reservation and makes the book available again
func (h *BookHandler) CancelReservation(c *gin.Context) {
	var req BookIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	if err := h.bookUC.CancelReservation(req.BookID, userUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Бронирование отменено"})
}

// GetShelfBooks returns books currently borrowed by the authenticated user
func (h *BookHandler) GetShelfBooks(c *gin.Context) {
	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	books, err := h.bookUC.GetBooksFromExchangesByStatus(userUUID, domain.ExchangeBorrowed)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, books)
}

// GetReadBooks returns books that the authenticated user has returned (read)
func (h *BookHandler) GetReadBooks(c *gin.Context) {
	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(401, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	books, err := h.bookUC.GetBooksFromExchangesByStatus(userUUID, domain.ExchangeReturned)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, books)
}

type ReturnBookRequest struct {
	Id                uuid.UUID            `json:"book_id" binding:"required"`
	Title             string               `json:"title" binding:"required"`
	Author            string               `json:"author" binding:"required"`
	Description       string               `json:"description"`
	Condition         domain.BookCondition `json:"condition" binding:"required,oneof=excellent good bad"`
	ImageURL          string               `json:"image_url"`
	CurrentLocationID *uuid.UUID           `json:"current_location_id"`
}

func (h *BookHandler) Return(c *gin.Context) {
	var req ReturnBookRequest

	// Similar tolerance as in Create: accept camelCase keys if present
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var bodyMap map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &bodyMap)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ImageURL == "" {
		if v, ok := bodyMap["imageUrl"].(string); ok && v != "" {
			req.ImageURL = v
		} else if v2, ok := bodyMap["image_url"].(string); ok && v2 != "" {
			req.ImageURL = v2
		}
	}

	if req.CurrentLocationID == nil {
		if v, ok := bodyMap["locationId"].(string); ok && v != "" {
			if parsed, err := uuid.Parse(v); err == nil {
				req.CurrentLocationID = &parsed
			}
		} else if v2, ok := bodyMap["current_location_id"].(string); ok && v2 != "" {
			if parsed, err := uuid.Parse(v2); err == nil {
				req.CurrentLocationID = &parsed
			}
		}
	}
	userIdRaw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}

	book := &domain.Book{
		ID:                req.Id,
		Title:             req.Title,
		Author:            req.Author,
		Description:       req.Description,
		Condition:         req.Condition,
		ImageURL:          req.ImageURL,
		CurrentLocationID: req.CurrentLocationID,
	}

	err := h.bookUC.Return(book, uuid.MustParse(userIDStr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно возвращена"})
}

type SetImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
}

// SetImage attaches an existing uploaded image (from ./data/images) to a book record.
// Accepts JSON { "image_url": "/images/<filename>" } or { "image_url": "<filename>" }.
func (h *BookHandler) SetImage(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	var req SetImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize filename
	filename := req.ImageURL
	if strings.HasPrefix(filename, "/images/") {
		filename = strings.TrimPrefix(filename, "/images/")
	}
	// prevent directory traversal
	filename = filepath.Base(filename)

	// check file exists in data/images
	imagesDir := "./data/images"
	fullPath := filepath.Join(imagesDir, filename)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Файл изображения не найден на сервере"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию о файле изображения"})
		return
	}

	// store public path
	public := "/images/" + filename
	if err := h.bookUC.SetBookImage(bookID, public); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Изображение прикреплено к книге", "image_url": public})
}

func (h *BookHandler) Borrow(c *gin.Context) {

	var req BookIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	err = h.bookUC.Borrow(req.BookID, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно выдана"})

}

func (h *BookHandler) Delete(c *gin.Context) {
	// Accept possible reason for deletion. We try to be tolerant with keys.
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req DeleteBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// tolerant parsing: accept book_id, bookId or id and reason as string
		var m map[string]interface{}
		if err2 := json.Unmarshal(bodyBytes, &m); err2 == nil {
			var idStr string
			if v, ok := m["book_id"].(string); ok && v != "" {
				idStr = v
			} else if v, ok := m["bookId"].(string); ok && v != "" {
				idStr = v
			} else if v, ok := m["id"].(string); ok && v != "" {
				idStr = v
			}
			if idStr != "" {
				if parsed, perr := uuid.Parse(idStr); perr == nil {
					req.BookID = parsed
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат идентификатора книги", "body": string(bodyBytes)})
					return
				}
			}
			if v, ok := m["reason"].(string); ok && v != "" {
				req.Reason = v
			} else if v2, ok := m["delete_reason"].(string); ok && v2 != "" {
				req.Reason = v2
			}
		}
		if req.BookID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	println("Book ID param:", req.BookID.String())

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя в токене"})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат идентификатора пользователя"})
		return
	}

	// determine if the caller is admin (middleware sets "role")
	isAdmin := false
	if roleRaw, exists := c.Get("role"); exists {
		if roleStr, ok := roleRaw.(string); ok && roleStr == string(domain.RoleAdmin) {
			isAdmin = true
		}
	}

	// pass reason (may be empty)
	err = h.bookUC.DeleteBook(req.BookID, userUUID, isAdmin, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно удалена"})

}
