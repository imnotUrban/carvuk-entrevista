package boleta

import (
	"errors"
	"fmt"
	"sort"

	"gorm.io/gorm"

	boletaentity "backend/entity/boleta"
	compraentity "backend/entity/compra"
	"backend/pkg/apperr"
	"backend/pkg/tax"
)

type ItemInput struct {
	ProductoID uint
	Cantidad   int
}

type CreateInput struct {
	Items []ItemInput
}

func (s *service) Create(input CreateInput) (*Detail, error) {
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("%w: items is required", apperr.ErrValidation)
	}
	if len(input.Items) > maxLines {
		return nil, fmt.Errorf("%w: at most %d items are allowed", apperr.ErrValidation, maxLines)
	}

	seen := make(map[uint]bool, len(input.Items))
	for _, it := range input.Items {
		if it.ProductoID == 0 {
			return nil, fmt.Errorf("%w: producto_id is required", apperr.ErrValidation)
		}
		if it.Cantidad < 1 || it.Cantidad > maxCantidad {
			return nil, fmt.Errorf("%w: cantidad must be between 1 and %d", apperr.ErrValidation, maxCantidad)
		}
		if seen[it.ProductoID] {
			return nil, fmt.Errorf("%w: duplicated producto_id %d", apperr.ErrValidation, it.ProductoID)
		}
		seen[it.ProductoID] = true
	}

	// Orden estable por producto para que las transacciones concurrentes bloqueen en el mismo orden.
	lines := append([]ItemInput(nil), input.Items...)
	sort.Slice(lines, func(i, j int) bool { return lines[i].ProductoID < lines[j].ProductoID })

	usuario, err := s.repo.FirstUsuario()
	if err != nil {
		return nil, err
	}

	compra := &compraentity.Compra{UserID: usuario.ID}
	bruto := 0
	for _, it := range lines {
		p, err := s.productos.GetByID(it.ProductoID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: producto %d", apperr.ErrNotFound, it.ProductoID)
			}
			return nil, err
		}
		if p.Stock < it.Cantidad {
			return nil, fmt.Errorf("%w: %q (available %d)", apperr.ErrInsufficientStock, p.Nombre, p.Stock)
		}
		compra.Items = append(compra.Items, compraentity.CompraItem{
			ProductoID:     p.ID,
			Nombre:         p.Nombre,
			PrecioUnitario: p.Precio,
			Cantidad:       it.Cantidad,
		})
		bruto += p.Precio * it.Cantidad
	}

	totals := tax.Calculate(bruto)
	boleta := &boletaentity.Boleta{
		ValorBruto:         totals.Bruto,
		Impuesto:           totals.Impuesto,
		ValorNeto:          totals.Neto,
		PorcentajeImpuesto: tax.TaxRatePercent,
	}
	if err := s.repo.Create(compra, boleta); err != nil {
		return nil, err
	}

	detail := &Detail{
		ID:                 boleta.ID,
		CompraID:           compra.ID,
		ValorBruto:         boleta.ValorBruto,
		Impuesto:           boleta.Impuesto,
		ValorNeto:          boleta.ValorNeto,
		PorcentajeImpuesto: boleta.PorcentajeImpuesto,
		Usuario:            usuario,
		CreatedAt:          boleta.CreatedAt,
	}
	for _, it := range compra.Items {
		detail.Items = append(detail.Items, DetailItem{
			ProductoID:     it.ProductoID,
			Nombre:         it.Nombre,
			PrecioUnitario: it.PrecioUnitario,
			Cantidad:       it.Cantidad,
			Subtotal:       it.PrecioUnitario * it.Cantidad,
		})
	}
	return detail, nil
}
