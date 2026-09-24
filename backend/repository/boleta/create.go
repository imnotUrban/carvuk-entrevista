package boleta

import (
	"fmt"

	"gorm.io/gorm"

	boletaentity "backend/entity/boleta"
	compraentity "backend/entity/compra"
	productoentity "backend/entity/producto"
	"backend/pkg/apperr"
)

func (r *repository) Create(compra *compraentity.Compra, boleta *boletaentity.Boleta) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Descuento atómico: solo afecta la fila si aún hay stock suficiente.
		for _, item := range compra.Items {
			res := tx.Model(&productoentity.Producto{}).
				Where("id = ? AND stock >= ?", item.ProductoID, item.Cantidad).
				Update("stock", gorm.Expr("stock - ?", item.Cantidad))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("%w: producto %d", apperr.ErrInsufficientStock, item.ProductoID)
			}
		}

		if err := tx.Create(compra).Error; err != nil {
			return err
		}

		boleta.CompraID = compra.ID
		return tx.Create(boleta).Error
	})
}
