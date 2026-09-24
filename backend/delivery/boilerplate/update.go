package boilerplate

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/boilerplate"
)

type updateRequest struct {
	Name     *string  `json:"name"`
	Code     *string  `json:"code"`
	Status   *string  `json:"status"`
	Quantity *int     `json:"quantity"`
	Amount   *float64 `json:"amount"`
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: "invalid id"})
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: err.Error()})
		return
	}

	item, err := h.service.Update(uint(id), service.UpdateInput{
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

	c.JSON(http.StatusOK, item)
}
