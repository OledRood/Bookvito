package http

import (
	"bookvito/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	// "strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookHandler struct {
	bookUC       domain.BookUseCase
	imageStorage domain.ImageStorage
}

const maxUploadImageSize = 20 * 1024 * 1024

func NewBookHandler(bookUC domain.BookUseCase, imageStorage domain.ImageStorage) *BookHandler {
	return &BookHandler{bookUC: bookUC, imageStorage: imageStorage}
}

// AutoFillBook returns metadata for a book using external provider (e.g., Google Books).
func (h *BookHandler) AutoFillBook(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Параметр q обязателен")
		return
	}

	meta, err := h.bookUC.AutoFill(query)
	if err != nil {
		respondBookUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, meta)
}

// UploadImage handles multipart image upload. Field name: "image".
// It saves file into ./data/images and returns JSON { "url": "/images/<filename>" }.
func (h *BookHandler) UploadImage(c *gin.Context) {
	if h.imageStorage == nil {
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, "Image storage is not configured")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Требуется файл изображения")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неподдерживаемый формат изображения")
		return
	}

	if file.Size <= 0 {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Файл изображения пустой")
		return
	}
	if file.Size > maxUploadImageSize {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Размер изображения не может превышать 20 MB")
		return
	}

	src, err := file.Open()
	if err != nil {
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, "Не удалось открыть файл изображения")
		return
	}
	defer src.Close()

	filename := uuid.New().String() + ext
	publicURL, err := h.imageStorage.Save(filename, file.Header.Get("Content-Type"), src)
	if err != nil {
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, "Не удалось сохранить файл")
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": publicURL})
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
		respondBookUseCaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetList(c *gin.Context) {
	filter, err := parseBookListFilter(c)
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	filter.OnlyAvailable = true
	if userIdRaw, exists := c.Get("userId"); exists {
		if userStr, ok := userIdRaw.(string); ok && userStr != "" {
			if parsed, err := uuid.Parse(userStr); err == nil {
				filter.ExcludeUserID = &parsed
			}
		}
	}

	response, err := h.bookUC.GetBooksList(filter)
	if err != nil {
		respondBookUseCaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func parseBookListFilter(c *gin.Context) (domain.BookListFilter, error) {
	filter := domain.BookListFilter{
		Search: strings.TrimSpace(c.Query("search")),
		SortBy: strings.ToLower(strings.TrimSpace(c.DefaultQuery("sort_by", "created_at"))),
		Order:  strings.ToLower(strings.TrimSpace(c.DefaultQuery("order", "desc"))),
		Limit:  20,
		Offset: 0,
	}

	if rawStatus := strings.ToLower(strings.TrimSpace(c.Query("status"))); rawStatus != "" {
		status := domain.BookStatus(rawStatus)
		switch status {
		case domain.BookAvailable, domain.BookRequested, domain.BookBorrowed, domain.BookArchived, domain.BookDeleted:
			filter.Status = &status
		default:
			return domain.BookListFilter{}, fmt.Errorf("неверный status: %s", rawStatus)
		}
	}

	if rawLocationID := strings.TrimSpace(c.Query("location_id")); rawLocationID != "" {
		locationID, err := uuid.Parse(rawLocationID)
		if err != nil {
			return domain.BookListFilter{}, fmt.Errorf("неверный location_id")
		}
		filter.LocationID = &locationID
	}

	allowedSortFields := map[string]struct{}{
		"title":               {},
		"author":              {},
		"status":              {},
		"created_at":          {},
		"updated_at":          {},
		"current_location_id": {},
	}
	if _, ok := allowedSortFields[filter.SortBy]; !ok {
		return domain.BookListFilter{}, fmt.Errorf("неверный sort_by: %s", filter.SortBy)
	}

	if filter.Order != "asc" && filter.Order != "desc" {
		return domain.BookListFilter{}, fmt.Errorf("неверный order: %s", filter.Order)
	}

	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil || limit <= 0 {
			return domain.BookListFilter{}, fmt.Errorf("limit должен быть положительным числом")
		}
		if limit > 100 {
			return domain.BookListFilter{}, fmt.Errorf("limit не может быть больше 100")
		}
		filter.Limit = limit
	}

	if rawOffset := strings.TrimSpace(c.Query("offset")); rawOffset != "" {
		offset, err := strconv.Atoi(rawOffset)
		if err != nil || offset < 0 {
			return domain.BookListFilter{}, fmt.Errorf("offset должен быть числом 0 или больше")
		}
		filter.Offset = offset
	}

	return filter, nil
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
		respondBookUseCaseError(c, err)
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
		respondBookUseCaseError(c, err)
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

	userID, isAdmin, err := getRequestUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.bookUC.GetBookStats(bookID, userID, isAdmin)
	if err != nil {
		respondBookUseCaseError(c, err)
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

type UpdateBookRequest struct {
	Title       *string               `json:"title"`
	Author      *string               `json:"author"`
	Description *string               `json:"description"`
	Condition   *domain.BookCondition `json:"condition"`
	ImageURL    *string               `json:"image_url"`
}

func (h *BookHandler) Update(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	userID, isAdmin, err := getRequestUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if len(bytes.TrimSpace(bodyBytes)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пустое тело запроса"})
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	forbiddenFields := map[string]struct{}{
		"id":         {},
		"owner_id":   {},
		"ownerId":    {},
		"owner":      {},
		"status":     {},
		"created_at": {},
		"createdAt":  {},
		"updated_at": {},
		"updatedAt":  {},
		"reviews":    {},
		"exchanges":  {},
	}
	allowedFields := map[string]struct{}{
		"title":               {},
		"author":              {},
		"description":         {},
		"condition":           {},
		"image_url":           {},
		"imageUrl":            {},
		"current_location_id": {},
		"currentLocationId":   {},
	}

	for key := range bodyMap {
		if _, forbidden := forbiddenFields[key]; forbidden {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Поле %s нельзя изменять", key)})
			return
		}
		if _, allowed := allowedFields[key]; !allowed {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Поле %s не поддерживается для обновления", key)})
			return
		}
	}

	if req.ImageURL == nil {
		if raw, ok := bodyMap["imageUrl"]; ok {
			var imageURL string
			if err := json.Unmarshal(raw, &imageURL); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный imageUrl"})
				return
			}
			req.ImageURL = &imageURL
		}
	}

	input := domain.BookUpdateInput{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Condition:   req.Condition,
		ImageURL:    req.ImageURL,
	}

	if raw, ok := bodyMap["current_location_id"]; ok {
		if err := applyCurrentLocation(raw, &input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if raw, ok := bodyMap["currentLocationId"]; ok {
		if err := applyCurrentLocation(raw, &input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if !hasBookUpdateFields(input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нужно передать хотя бы одно поле для обновления"})
		return
	}

	book, err := h.bookUC.UpdateBook(bookID, userID, isAdmin, input)
	if err != nil {
		respondBookUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func applyCurrentLocation(raw json.RawMessage, input *domain.BookUpdateInput) error {
	input.CurrentLocationIDSet = true
	if string(raw) == "null" {
		input.CurrentLocationID = nil
		return nil
	}

	var locationID string
	if err := json.Unmarshal(raw, &locationID); err != nil {
		return fmt.Errorf("неверный current_location_id")
	}

	locationUUID, err := uuid.Parse(strings.TrimSpace(locationID))
	if err != nil {
		return fmt.Errorf("неверный current_location_id")
	}

	input.CurrentLocationID = &locationUUID
	return nil
}

func hasBookUpdateFields(input domain.BookUpdateInput) bool {
	return input.Title != nil ||
		input.Author != nil ||
		input.Description != nil ||
		input.Condition != nil ||
		input.ImageURL != nil ||
		input.CurrentLocationIDSet
}

func getRequestUser(c *gin.Context) (uuid.UUID, bool, error) {
	userIdRaw, ok := c.Get("userId")
	if !ok {
		return uuid.Nil, false, domain.NewUnauthorizedError("Пользователь не аутентифицирован")
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, false, domain.NewUnauthorizedError("Неверный идентификатор пользователя в токене")
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, false, domain.NewUnauthorizedError("Неверный формат идентификатора пользователя")
	}

	isAdmin := false
	if roleRaw, exists := c.Get("role"); exists {
		if roleStr, ok := roleRaw.(string); ok && roleStr == string(domain.RoleAdmin) {
			isAdmin = true
		}
	}

	return userUUID, isAdmin, nil
}

func respondBookUseCaseError(c *gin.Context, err error) {
	WriteErrorFromErr(c, err)
}

func extractOptionalLocationID(bodyMap map[string]json.RawMessage, keys ...string) (*uuid.UUID, error) {
	for _, key := range keys {
		raw, ok := bodyMap[key]
		if !ok {
			continue
		}
		if string(raw) == "null" {
			return nil, nil
		}

		var locationID string
		if err := json.Unmarshal(raw, &locationID); err != nil {
			return nil, fmt.Errorf("неверный current_location_id")
		}

		locationUUID, err := uuid.Parse(strings.TrimSpace(locationID))
		if err != nil {
			return nil, fmt.Errorf("неверный current_location_id")
		}
		return &locationUUID, nil
	}

	return nil, nil
}

func (h *BookHandler) Create(c *gin.Context) {
	var req CreateBookRequest

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Некорректный JSON")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	if req.ImageURL == "" {
		if raw, ok := bodyMap["imageUrl"]; ok {
			if err := json.Unmarshal(raw, &req.ImageURL); err != nil {
				WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный imageUrl")
				return
			}
		} else if raw, ok := bodyMap["image_url"]; ok {
			if err := json.Unmarshal(raw, &req.ImageURL); err != nil {
				WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный image_url")
				return
			}
		}
	}

	if req.CurrentLocationID == nil {
		locationID, err := extractOptionalLocationID(bodyMap, "locationId", "current_location_id", "currentLocationId")
		if err != nil {
			WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
			return
		}
		req.CurrentLocationID = locationID
	}

	userIdRaw, exists := c.Get("userId")
	if !exists {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Неверный идентификатор пользователя в токене")
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
		respondBookUseCaseError(c, err)
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
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный location_id"})
			return
		}
	} else if v2, ok := bodyMap["location_id"].(string); ok && v2 != "" {
		if parsed, perr := uuid.Parse(v2); perr == nil {
			locIDPtr = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный location_id"})
			return
		}
	}

	err = h.bookUC.Request(req.BookID, userUUID, locIDPtr)
	if err != nil {
		respondBookUseCaseError(c, err)
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
		respondBookUseCaseError(c, err)
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
		respondBookUseCaseError(c, err)
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
	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Некорректный JSON")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	if req.ImageURL == "" {
		if raw, ok := bodyMap["imageUrl"]; ok {
			if err := json.Unmarshal(raw, &req.ImageURL); err != nil {
				WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный imageUrl")
				return
			}
		} else if raw, ok := bodyMap["image_url"]; ok {
			if err := json.Unmarshal(raw, &req.ImageURL); err != nil {
				WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный image_url")
				return
			}
		}
	}

	if req.CurrentLocationID == nil {
		locationID, err := extractOptionalLocationID(bodyMap, "locationId", "current_location_id", "currentLocationId")
		if err != nil {
			WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
			return
		}
		req.CurrentLocationID = locationID
	}
	userID, isAdmin, err := getRequestUser(c)
	if err != nil {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, err.Error())
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

	err = h.bookUC.Return(book, userID, isAdmin)
	if err != nil {
		respondBookUseCaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно возвращена"})
}

type SetImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
}

func (h *BookHandler) SetImage(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный идентификатор книги")
		return
	}

	var req SetImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	userID, isAdmin, err := getRequestUser(c)
	if err != nil {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, err.Error())
		return
	}

	if err := h.bookUC.SetBookImage(bookID, userID, isAdmin, req.ImageURL); err != nil {
		respondBookUseCaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Изображение прикреплено к книге", "image_url": req.ImageURL})
}

func (h *BookHandler) DeleteImage(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный идентификатор книги")
		return
	}

	userID, isAdmin, err := getRequestUser(c)
	if err != nil {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, err.Error())
		return
	}

	if err := h.bookUC.DeleteBookImage(bookID, userID, isAdmin); err != nil {
		respondBookUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение книги удалено"})
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
		respondBookUseCaseError(c, err)
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
		respondBookUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно удалена"})

}
