package boleta

import (
	"fmt"
	"strings"
	"time"

	usuarioentity "backend/entity/usuario"
	"backend/pkg/apperr"
	repo "backend/repository/boleta"
)

type ListInput struct {
	Page  int
	Limit int
	Order string
}

// ListItem es el resumen de una boleta en el listado.
type ListItem struct {
	ID            uint                   `json:"id"`
	CreatedAt     time.Time              `json:"created_at"`
	ValorBruto    int                    `json:"valor_bruto"`
	CantidadItems int                    `json:"cantidad_items"`
	Usuario       *usuarioentity.Usuario `json:"usuario"`
}

func (s *service) List(input ListInput) ([]ListItem, int64, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	order := strings.ToLower(input.Order)
	if order == "" {
		order = "desc"
	}
	if order != "asc" && order != "desc" {
		return nil, 0, fmt.Errorf("%w: order must be asc or desc", apperr.ErrValidation)
	}

	boletas, total, err := s.repo.List(repo.Filter{Page: page, Limit: limit, Order: order})
	if err != nil {
		return nil, 0, err
	}

	items := make([]ListItem, 0, len(boletas))
	for _, b := range boletas {
		li := ListItem{ID: b.ID, CreatedAt: b.CreatedAt, ValorBruto: b.ValorBruto}
		if b.Compra != nil {
			li.Usuario = b.Compra.Usuario
			for _, it := range b.Compra.Items {
				li.CantidadItems += it.Cantidad
			}
		}
		items = append(items, li)
	}
	return items, total, nil
}
