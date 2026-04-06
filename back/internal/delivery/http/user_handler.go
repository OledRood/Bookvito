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
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Неверный идентификатор пользователя в токене")
		return
	}
	user, err := h.userUC.GetUserByID(userIDStr)
	if err != nil {
		WriteErrorFromErr(c, err)
		return
	}
	if user == nil {
		WriteError(c, http.StatusNotFound, domain.ErrorCodeNotFound, "Пользователь не найден")
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
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	tokens, err := h.userUC.RegisterUser(req.Email, req.Password, req.Name)
	if err != nil {
		WriteErrorFromErr(c, err)
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
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}
	tokens, err := h.userUC.LoginUser(req.Email, req.Password)
	if err != nil {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, err.Error())
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
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	tokens, err := h.userUC.RefreshToken(req.RefreshToken)
	if err != nil {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Идентификатор пользователя не найден")
		return
	}
	if err := h.userUC.Logout(userID.(string)); err != nil {
		WriteErrorFromErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Сессия завершена"})
}

func (h *UserHandler) GetMyMovementHistory(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Идентификатор пользователя не найден в контексте")
		return
	}

	history, err := h.userUC.GetUserMovementHistory(userID.(string))
	if err != nil {
		WriteErrorFromErr(c, err)
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
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	userIDStr, ok := userIdRaw.(string)
	if !ok || userIDStr == "" {
		WriteError(c, http.StatusUnauthorized, domain.ErrorCodeUnauthorized, "Неверный идентификатор пользователя в токене")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	user, err := h.userUC.GetUserByID(userIDStr)
	if err != nil || user == nil {
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, "Не удалось загрузить пользователя")
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
		WriteErrorFromErr(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteMe deletes the currently authenticated user
func (h *UserHandler) DeleteMe(c *gin.Context) {
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

	if err := h.userUC.DeleteUser(userIDStr); err != nil {
		WriteErrorFromErr(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
