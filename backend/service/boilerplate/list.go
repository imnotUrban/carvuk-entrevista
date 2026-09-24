package boilerplate

import (
	"strings"

	entity "backend/entity/boilerplate"
	repo "backend/repository/boilerplate"
)

type ListInput struct {
	Name      string
	Code      string
	Status    string
	MinAmount *float64
	MaxAmount *float64
	Page      int
	Limit     int
	SortBy    string
	Order     string
}

func (s *service) List(input ListInput) ([]entity.Boilerplate, int64, error) {
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
	if !sortWhitelist[sortBy] {
		sortBy = "id"
	}
	order := strings.ToLower(input.Order)
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	status := input.Status
	if !statusWhitelist[status] {
		status = ""
	}

	filter := repo.Filter{
		Name:      strings.TrimSpace(input.Name),
		Code:      strings.TrimSpace(input.Code),
		Status:    status,
		MinAmount: input.MinAmount,
		MaxAmount: input.MaxAmount,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		Order:     order,
	}

	return s.repo.List(filter)
}
