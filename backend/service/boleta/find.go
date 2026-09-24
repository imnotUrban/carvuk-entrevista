package boleta

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"backend/pkg/apperr"
)

func (s *service) GetByID(id uint) (*Detail, error) {
	b, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: boleta %d", apperr.ErrNotFound, id)
		}
		return nil, err
	}

	detail := &Detail{
		ID:                 b.ID,
		CompraID:           b.CompraID,
		ValorBruto:         b.ValorBruto,
		Impuesto:           b.Impuesto,
		ValorNeto:          b.ValorNeto,
		PorcentajeImpuesto: b.PorcentajeImpuesto,
		CreatedAt:          b.CreatedAt,
		Items:              []DetailItem{},
	}
	if b.Compra != nil {
		detail.Usuario = b.Compra.Usuario
		for _, it := range b.Compra.Items {
			detail.Items = append(detail.Items, DetailItem{
				ProductoID:     it.ProductoID,
				Nombre:         it.Nombre,
				PrecioUnitario: it.PrecioUnitario,
				Cantidad:       it.Cantidad,
				Subtotal:       it.PrecioUnitario * it.Cantidad,
			})
		}
	}
	return detail, nil
}
