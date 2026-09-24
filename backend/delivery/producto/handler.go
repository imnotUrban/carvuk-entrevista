package producto

import (
	"github.com/gin-gonic/gin"

	service "backend/service/producto"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	items := rg.Group("/productos")
	{
		items.POST("", h.Create)
		items.GET("", h.List)
		items.GET("/:id", h.Get)
		items.PUT("/:id", h.Update)
		items.DELETE("/:id", h.Delete)
	}
}
