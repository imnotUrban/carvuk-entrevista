package producto

import (
	"fmt"
	"strings"

	entity "backend/entity/producto"
	"backend/pkg/apperr"
)

type UpdateInput struct {
	Nombre *string
	Precio *int
	Stock  *int
}

func (s *service) Update(id uint, input UpdateInput) (*entity.Producto, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Nombre != nil {
		nombre := strings.TrimSpace(*input.Nombre)
		if nombre == "" {
			return nil, fmt.Errorf("%w: nombre cannot be empty", apperr.ErrValidation)
		}
		item.Nombre = nombre
	}
	if input.Precio != nil {
		if err := validatePrecio(*input.Precio); err != nil {
			return nil, err
		}
		item.Precio = *input.Precio
	}
	if input.Stock != nil {
		if err := validateStock(*input.Stock); err != nil {
			return nil, err
		}
		item.Stock = *input.Stock
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(item.ID)
}
