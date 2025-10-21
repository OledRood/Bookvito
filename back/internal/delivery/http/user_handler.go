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
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
		return
	}
	user, err := h.userUC.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
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
	Name     string `json:"name" binding:"required,min=5,max=50"`
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

func (h *UserHandler) GetMyMovementHistory(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
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

// // GetByID retrieves a user by ID

// func (h *UserHandler) GetByID(c *gin.Context) {
// 	userID, ok := c.Get("userId")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
// 		return
// 	}
// 	user, err := h.userUC.GetUserByID(userID.(string))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// Update updates a user
//
//	func (h *UserHandler) Update(c *gin.Context) {
//		userID, ok := c.Get("userId")
// func (h *UserHandler) Update(c *gin.Context) {
// 	userID, ok := c.Get("userId")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
// 		return
// 	}
// 	user, err := h.userUC.GetUserByID(userID.(string))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}
// 	if err := c.ShouldBindJSON(user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := h.userUC.UpdateUser(user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
// 		return
// 	}
// 	user, err := h.userUC.GetUserByID(userID.(string))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}
// 	if err := c.ShouldBindJSON(user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := h.userUC.UpdateUser(user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// // Delete deletes a user
//
//	func (h *UserHandler) Delete(c *gin.Context) {
//		userID, ok := c.Get("userId")
// func (h *UserHandler) Delete(c *gin.Context) {
// 	userID, ok := c.Get("userId")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
// 		return
// 	}
// 	if err := h.userUC.DeleteUser(userID.(string)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
// }

// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
// 		return
// 	}
// 	if err := h.userUC.DeleteUser(userID.(string)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
// }

// // List retrieves a list of users
// func (h *UserHandler) List(c *gin.Context) {
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

// 	users, err := h.userUC.ListUsers(limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, users)
// }
