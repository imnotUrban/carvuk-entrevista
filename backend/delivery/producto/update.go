package producto

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/producto"
)

type updateRequest struct {
	Nombre *string `json:"nombre"`
	Precio *int    `json:"precio"`
	Stock  *int    `json:"stock"`
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
		Nombre: req.Nombre,
		Precio: req.Precio,
		Stock:  req.Stock,
	})
	if err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}
