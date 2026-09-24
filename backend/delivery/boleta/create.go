package boleta

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/boleta"
)

// Solo se aceptan producto_id y cantidad: precios, impuesto y totales los calcula el backend.
type itemRequest struct {
	ProductoID uint `json:"producto_id"`
	Cantidad   int  `json:"cantidad"`
}

type createRequest struct {
	Items []itemRequest `json:"items"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: err.Error()})
		return
	}

	input := service.CreateInput{}
	for _, it := range req.Items {
		input.Items = append(input.Items, service.ItemInput{ProductoID: it.ProductoID, Cantidad: it.Cantidad})
	}

	detail, err := h.service.Create(input)
	if err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, detail)
}
