package producto

import (
	"fmt"
	"strings"

	entity "backend/entity/producto"
	"backend/pkg/apperr"
)

type CreateInput struct {
	Nombre string
	Precio int
	Stock  int
}

func (s *service) Create(input CreateInput) (*entity.Producto, error) {
	nombre := strings.TrimSpace(input.Nombre)
	if nombre == "" {
		return nil, fmt.Errorf("%w: nombre is required", apperr.ErrValidation)
	}
	if err := validatePrecio(input.Precio); err != nil {
		return nil, err
	}
	if err := validateStock(input.Stock); err != nil {
		return nil, err
	}

	item := &entity.Producto{Nombre: nombre, Precio: input.Precio, Stock: input.Stock}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(item.ID)
}
