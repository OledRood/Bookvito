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
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}
	if !checkAdminRole(c) {
		WriteError(c, http.StatusForbidden, domain.ErrorCodeForbidden, "only admin can create locations")
		return
	}

	location := &domain.Location{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := h.locationUC.Create(location); err != nil {
		WriteErrorFromErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "location created successfully"})
}

func (h *LocationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "invalid location ID")
		return
	}

	location, err := h.locationUC.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			WriteError(c, http.StatusNotFound, domain.ErrorCodeNotFound, "location not found")
			return
		}
		WriteErrorFromErr(c, err)
		return
	}
	if location == nil {
		WriteError(c, http.StatusNotFound, domain.ErrorCodeNotFound, "location not found")
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) GetAll(c *gin.Context) {
	locations, err := h.locationUC.GetAll()
	if err != nil {
		WriteErrorFromErr(c, err)
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "invalid location ID")
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
	}
	if !checkAdminRole(c) {
		WriteError(c, http.StatusForbidden, domain.ErrorCodeForbidden, "only admin can create locations")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, err.Error())
		return
	}

	location := &domain.Location{
		ID:      id,
		Name:    req.Name,
		Address: req.Address,
	}

	if err := h.locationUC.Update(location); err != nil {
		WriteErrorFromErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location updated successfully"})
}

func (h *LocationHandler) Delete(c *gin.Context) {

	if !checkAdminRole(c) {
		WriteError(c, http.StatusForbidden, domain.ErrorCodeForbidden, "only admin can create locations")
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		WriteError(c, http.StatusBadRequest, domain.ErrorCodeValidation, "invalid location ID")
		return
	}

	if err := h.locationUC.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			WriteError(c, http.StatusNotFound, domain.ErrorCodeNotFound, "location not found")
			return
		}
		WriteErrorFromErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location deleted successfully"})
}

func checkAdminRole(c *gin.Context) bool {
	v, ok := c.Get("role")
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	return s == "admin"
}
