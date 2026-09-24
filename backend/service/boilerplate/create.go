package boilerplate

import (
	"fmt"
	"strings"

	entity "backend/entity/boilerplate"
	"backend/pkg/apperr"
)

type CreateInput struct {
	Name     string
	Code     string
	Status   string
	Quantity int
	Amount   float64
}

func (s *service) Create(input CreateInput) (*entity.Boilerplate, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", apperr.ErrValidation)
	}
	code := strings.TrimSpace(input.Code)
	if code == "" {
		return nil, fmt.Errorf("%w: code is required", apperr.ErrValidation)
	}
	status := input.Status
	if status == "" {
		status = entity.StatusDraft
	}
	if err := validateStatus(status); err != nil {
		return nil, err
	}
	if input.Quantity < 0 {
		return nil, fmt.Errorf("%w: quantity cannot be negative", apperr.ErrValidation)
	}
	if input.Amount < 0 {
		return nil, fmt.Errorf("%w: amount cannot be negative", apperr.ErrValidation)
	}

	duplicate, err := s.repo.ExistsByCode(code, 0)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, apperr.ErrDuplicateCode
	}

	item := &entity.Boilerplate{
		Name:     name,
		Code:     code,
		Status:   status,
		Quantity: input.Quantity,
		Amount:   input.Amount,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(item.ID)
}
