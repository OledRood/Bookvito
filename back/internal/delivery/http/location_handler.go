package http

import (
	"bookvito/internal/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationHandler struct {
	locationUC domain.LocationUseCase
}

func NewLocationHandler(locationUC domain.LocationUseCase) *LocationHandler {
	return &LocationHandler{locationUC: locationUC}
}

func (h *LocationHandler) Create(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !checkAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can create locations"})
		return
	}

	location := &domain.Location{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := h.locationUC.Create(location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "location created successfully"})
}

func (h *LocationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	location, err := h.locationUC.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if location == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) GetAll(c *gin.Context) {
	locations, err := h.locationUC.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
	}
	if !checkAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can create locations"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location := &domain.Location{
		ID:      id,
		Name:    req.Name,
		Address: req.Address,
	}

	if err := h.locationUC.Update(location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location updated successfully"})
}

func (h *LocationHandler) Delete(c *gin.Context) {

	if !checkAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can create locations"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	if err := h.locationUC.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location deleted successfully"})
}

func checkAdminRole(c *gin.Context) bool {
	userRole, exists := c.Get("role")
	if !exists {
		return false
	}
	return userRole == "admin"
}
