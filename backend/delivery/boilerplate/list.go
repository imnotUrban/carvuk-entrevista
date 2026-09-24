package boilerplate

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
	service "backend/service/boilerplate"
)

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	var minAmount, maxAmount *float64
	if v := c.Query("min_amount"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minAmount = &parsed
		}
	}
	if v := c.Query("max_amount"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxAmount = &parsed
		}
	}

	items, total, err := h.service.List(service.ListInput{
		Name:      c.Query("name"),
		Code:      c.Query("code"),
		Status:    c.Query("status"),
		MinAmount: minAmount,
		MaxAmount: maxAmount,
		Page:      page,
		Limit:     limit,
		SortBy:    c.Query("sort_by"),
		Order:     c.Query("order"),
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

	c.JSON(http.StatusOK, httpkit.PaginatedResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
