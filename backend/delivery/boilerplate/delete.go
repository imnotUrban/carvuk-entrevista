package boilerplate

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/pkg/httpkit"
)

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpkit.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		httpkit.HandleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
