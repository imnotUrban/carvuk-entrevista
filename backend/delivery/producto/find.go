package producto

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
)

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: "invalid id"})
		return
	}

	item, err := h.service.GetByID(uint(id))
	if err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}
