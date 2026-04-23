package http

import (
	"bookvito/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExchangeHandler struct {
	exchangeUC *usecase.ExchangeUseCase
}

// NewExchangeHandler creates a new exchange handler
func NewExchangeHandler(exchangeUC *usecase.ExchangeUseCase) *ExchangeHandler {
	return &ExchangeHandler{exchangeUC: exchangeUC}
}

type CreateExchangeRequest struct {
	RequesterID uint   `json:"requester_id" binding:"required"`
	BookID      uint   `json:"book_id" binding:"required"`
	Message     string `json:"message"`
}

// Create creates a new exchange request
// FIX: Uncommented to enable exchange endpoints for tests
func (h *ExchangeHandler) Create(c *gin.Context) {
	var req CreateExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchange, err := h.exchangeUC.CreateExchangeRequest(req.RequesterID, req.BookID, req.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exchange)
}

// GetByID retrieves an exchange by ID
// FIX: Uncommented to enable exchange endpoints for tests
func (h *ExchangeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exchange ID"})
		return
	}

	exchange, err := h.exchangeUC.GetExchangeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exchange not found"})
		return
	}

	c.JSON(http.StatusOK, exchange)
}

// List retrieves a list of exchanges
// FIX: Uncommented to enable exchange endpoints for tests
func (h *ExchangeHandler) List(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	exchanges, err := h.exchangeUC.GetExchangesByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exchanges)
}

type UpdateExchangeStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus updates the status of an exchange
// FIX: Uncommented to enable exchange endpoints for tests
func (h *ExchangeHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exchange ID"})
		return
	}

	var req UpdateExchangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchange, err := h.exchangeUC.UpdateExchangeStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exchange)
}

// Delete deletes an exchange
// FIX: Uncommented to enable exchange endpoints for tests
func (h *ExchangeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exchange ID"})
		return
	}

	if err := h.exchangeUC.DeleteExchange(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
