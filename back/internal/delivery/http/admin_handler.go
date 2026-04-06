package http

import (
	"bookvito/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userUC  domain.UserUseCase
	adminUC domain.AdminUseCase
}

func NewAdminHandler(userUC domain.UserUseCase, adminUC domain.AdminUseCase) *AdminHandler {
	return &AdminHandler{userUC: userUC, adminUC: adminUC}
}

// GetStats GET /admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminUC.GetStats()
	if err != nil {
		WriteErrorFromErr(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ListUsers GET /admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.userUC.ListUsers(100, 0)
	if err != nil {
		WriteErrorFromErr(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUserRole PUT /admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Неверный идентификатор пользователя")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	role := domain.UserRole(req.Role)
	if role != domain.RoleUser && role != domain.RoleModer && role != domain.RoleAdmin {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "Недопустимое значение роли. Допустимые: user, moder, admin")
		return
	}

	if err := h.userUC.UpdateUserRole(userID, role); err != nil {
		WriteErrorFromErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Роль пользователя обновлена"})
}
