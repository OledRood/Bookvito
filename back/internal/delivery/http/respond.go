package http

import (
	"bookvito/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string           `json:"error"`
	Code  domain.ErrorCode `json:"code"`
}

func WriteError(c *gin.Context, status int, code domain.ErrorCode, message string) {
	if code == "" {
		code = domain.ErrorCodeInternal
	}
	c.JSON(status, errorResponse{
		Error: message,
		Code:  code,
	})
}

func WriteErrorFromErr(c *gin.Context, err error) {
	if err == nil {
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, "internal error")
		return
	}

	code, ok := domain.AppErrorCode(err)
	if !ok {
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, err.Error())
		return
	}

	switch code {
	case domain.ErrorCodeValidation:
		WriteError(c, http.StatusBadRequest, code, err.Error())
	case domain.ErrorCodeUnauthorized:
		WriteError(c, http.StatusUnauthorized, code, err.Error())
	case domain.ErrorCodeForbidden:
		WriteError(c, http.StatusForbidden, code, err.Error())
	case domain.ErrorCodeNotFound:
		WriteError(c, http.StatusNotFound, code, err.Error())
	case domain.ErrorCodeConflict:
		WriteError(c, http.StatusConflict, code, err.Error())
	default:
		WriteError(c, http.StatusInternalServerError, domain.ErrorCodeInternal, err.Error())
	}
}
