package httpkit

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/pkg/apperr"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

func HandleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apperr.ErrValidation):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apperr.ErrDuplicateCode):
		c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
