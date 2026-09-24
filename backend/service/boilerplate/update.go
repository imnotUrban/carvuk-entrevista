package boilerplate

import (
	"fmt"
	"strings"

	entity "backend/entity/boilerplate"
	"backend/pkg/apperr"
)

type UpdateInput struct {
	Name     *string
	Code     *string
	Status   *string
	Quantity *int
	Amount   *float64
}

func (s *service) Update(id uint, input UpdateInput) (*entity.Boilerplate, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, fmt.Errorf("%w: name cannot be empty", apperr.ErrValidation)
		}
		item.Name = name
	}
	if input.Status != nil {
		if err := validateStatus(*input.Status); err != nil {
			return nil, err
		}
		item.Status = *input.Status
	}
	if input.Quantity != nil {
		if *input.Quantity < 0 {
			return nil, fmt.Errorf("%w: quantity cannot be negative", apperr.ErrValidation)
		}
		item.Quantity = *input.Quantity
	}
	if input.Amount != nil {
		if *input.Amount < 0 {
			return nil, fmt.Errorf("%w: amount cannot be negative", apperr.ErrValidation)
		}
		item.Amount = *input.Amount
	}
	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return nil, fmt.Errorf("%w: code cannot be empty", apperr.ErrValidation)
		}
		duplicate, err := s.repo.ExistsByCode(code, item.ID)
		if err != nil {
			return nil, err
		}
		if duplicate {
			return nil, apperr.ErrDuplicateCode
		}
		item.Code = code
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(item.ID)
}
