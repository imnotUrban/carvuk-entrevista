package producto

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/producto"
)

type createRequest struct {
	Nombre string `json:"nombre" binding:"required"`
	Precio int    `json:"precio"`
	Stock  int    `json:"stock"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: err.Error()})
		return
	}

	item, err := h.service.Create(service.CreateInput{
		Nombre: req.Nombre,
		Precio: req.Precio,
		Stock:  req.Stock,
	})
	if err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}
