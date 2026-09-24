package producto

import (
	"fmt"
	"strings"

	entity "backend/entity/producto"
	"backend/pkg/apperr"
	repo "backend/repository/producto"
)

type ListInput struct {
	Nombre string
	Page   int
	Limit  int
	SortBy string
	Order  string
}

func (s *service) List(input ListInput) ([]entity.Producto, int64, error) {
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

	sortBy := input.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	if !sortWhitelist[sortBy] {
		return nil, 0, fmt.Errorf("%w: invalid sort_by", apperr.ErrValidation)
	}
	order := strings.ToLower(input.Order)
	if order == "" {
		order = "asc"
	}
	if order != "asc" && order != "desc" {
		return nil, 0, fmt.Errorf("%w: order must be asc or desc", apperr.ErrValidation)
	}

	return s.repo.List(repo.Filter{
		Nombre: strings.TrimSpace(input.Nombre),
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  order,
	})
}
