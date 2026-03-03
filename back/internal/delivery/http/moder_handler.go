package http

import (
	"bookvito/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ModerHandler struct {
	moderUC domain.ModerUseCase
}

func NewModerHandler(moderUC domain.ModerUseCase) *ModerHandler {
	return &ModerHandler{moderUC: moderUC}
}

// GetReports GET /moder/reports
func (h *ModerHandler) GetReports(c *gin.Context) {
	reports, err := h.moderUC.GetReports(100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// CreateReport POST /books/:id/report
func (h *ModerHandler) CreateReport(c *gin.Context) {
	idParam := c.Param("id")
	bookID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	userIdRaw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	userID, err := uuid.Parse(userIdRaw.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.moderUC.CreateReport(bookID, userID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Жалоба отправлена"})
}

// ResolveReport PUT /moder/reports/:id/resolve
func (h *ModerHandler) ResolveReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор жалобы"})
		return
	}
	if err := h.moderUC.ResolveReport(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Жалоба закрыта"})
}

// DismissReport PUT /moder/reports/:id/dismiss
func (h *ModerHandler) DismissReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор жалобы"})
		return
	}
	if err := h.moderUC.DismissReport(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Жалоба отклонена"})
}

// ArchiveBook PUT /moder/books/:id/archive
func (h *ModerHandler) ArchiveBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}
	if err := h.moderUC.ArchiveBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Книга архивирована"})
}
