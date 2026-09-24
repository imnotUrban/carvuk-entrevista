package producto

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/producto"
)

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	items, total, err := h.service.List(service.ListInput{
		Nombre: c.Query("nombre"),
		Page:   page,
		Limit:  limit,
		SortBy: c.Query("sort_by"),
		Order:  c.Query("order"),
	})
	if err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	c.JSON(http.StatusOK, httpkit.PaginatedResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
