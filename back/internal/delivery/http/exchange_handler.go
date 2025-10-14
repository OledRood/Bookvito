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
func (h *ExchangeHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	exchanges, err := h.exchangeUC.ListExchanges(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exchanges)
}

// GetByRequester retrieves exchanges by requester
func (h *ExchangeHandler) GetByRequester(c *gin.Context) {
	requesterID, err := strconv.ParseUint(c.Param("requester_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester ID"})
		return
	}

	exchanges, err := h.exchangeUC.GetExchangesByRequester(uint(requesterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exchanges)
}

// GetByOwner retrieves exchanges by owner
func (h *ExchangeHandler) GetByOwner(c *gin.Context) {
	ownerID, err := strconv.ParseUint(c.Param("owner_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner ID"})
		return
	}

	exchanges, err := h.exchangeUC.GetExchangesByOwner(uint(ownerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exchanges)
}

// Accept accepts an exchange request
func (h *ExchangeHandler) Accept(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exchange ID"})
		return
	}

	ownerID, _ := strconv.ParseUint(c.Query("owner_id"), 10, 32)

	if err := h.exchangeUC.AcceptExchange(uint(id), uint(ownerID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exchange accepted successfully"})
}

// Reject rejects an exchange request
func (h *ExchangeHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exchange ID"})
		return
	}

	ownerID, _ := strconv.ParseUint(c.Query("owner_id"), 10, 32)

	if err := h.exchangeUC.RejectExchange(uint(id), uint(ownerID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exchange rejected successfully"})
}

// Complete completes an exchange
func (h *ExchangeHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exchange ID"})
		return
	}

	ownerID, _ := strconv.ParseUint(c.Query("owner_id"), 10, 32)

	if err := h.exchangeUC.CompleteExchange(uint(id), uint(ownerID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exchange completed successfully"})
}
