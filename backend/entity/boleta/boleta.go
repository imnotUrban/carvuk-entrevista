package boleta

import (
	"time"

	compraentity "backend/entity/compra"
)

type Boleta struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	CompraID           uint      `gorm:"column:id_compra;not null;index" json:"compra_id"`
	ValorNeto          int       `gorm:"not null" json:"valor_neto"`
	ValorBruto         int       `gorm:"not null" json:"valor_bruto"`
	Impuesto           int       `gorm:"not null" json:"impuesto"`
	PorcentajeImpuesto int       `gorm:"not null" json:"porcentaje_impuesto"`
	CreatedAt          time.Time `json:"created_at"`

	Compra *compraentity.Compra `gorm:"foreignKey:CompraID" json:"-"`
}

func (Boleta) TableName() string {
	return "boletas"
}
