package boilerplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/boilerplate"
)

type createRequest struct {
	Name     string  `json:"name" binding:"required"`
	Code     string  `json:"code" binding:"required"`
	Status   string  `json:"status"`
	Quantity int     `json:"quantity"`
	Amount   float64 `json:"amount"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: err.Error()})
		return
	}

	item, err := h.service.Create(service.CreateInput{
		Name:     req.Name,
		Code:     req.Code,
		Status:   req.Status,
		Quantity: req.Quantity,
		Amount:   req.Amount,
	})
	if err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}
