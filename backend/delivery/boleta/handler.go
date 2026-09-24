package boleta

import (
	"github.com/gin-gonic/gin"

	service "backend/service/boleta"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	items := rg.Group("/boletas")
	{
		items.POST("", h.Create)
		items.GET("", h.List)
		items.GET("/:id", h.Get)
	}
}
