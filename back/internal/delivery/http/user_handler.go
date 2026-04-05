package http

import (
	"bookvito/internal/domain"
	"net/http"

	// "strconv"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC domain.UserUseCase
}

func (h *UserHandler) GetByID(c *gin.Context) {

	println("GetByID called")
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
	user, err := h.userUC.GetUserByID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// фабрика создания нового обработчика пользователей
func NewUserHandler(userUC domain.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=2,max=50"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.userUC.RegisterUser(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.userUC.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.userUC.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Идентификатор пользователя не найден"})
		return
	}
	if err := h.userUC.Logout(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Сессия завершена"})
}

func (h *UserHandler) GetMyMovementHistory(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Идентификатор пользователя не найден в контексте"})
		return
	}

	history, err := h.userUC.GetUserMovementHistory(userID.(string))
	if err != nil {
		// В usecase уже есть проверка на формат UUID, но на всякий случай
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

type UpdateUserRequest struct {
	Avatar *string `json:"avatar"`
	Name   *string `json:"name"`
}

// UpdateMe updates fields of the currently authenticated user (partial update supported)
func (h *UserHandler) UpdateMe(c *gin.Context) {
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

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUC.GetUserByID(userIDStr)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить пользователя"})
		return
	}

	// Only update allowed fields (avatar for now)
	// Update allowed fields
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Name != nil {
		user.Name = *req.Name
	}

	if err := h.userUC.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteMe deletes the currently authenticated user
func (h *UserHandler) DeleteMe(c *gin.Context) {
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

	if err := h.userUC.DeleteUser(userIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
