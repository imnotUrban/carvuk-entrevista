package compra

import (
	"time"

	usuarioentity "backend/entity/usuario"
)

type Compra struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	UserID    uint                   `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time              `json:"created_at"`
	Items     []CompraItem           `gorm:"foreignKey:CompraID" json:"items"`
	Usuario   *usuarioentity.Usuario `gorm:"foreignKey:UserID" json:"-"`
}

func (Compra) TableName() string {
	return "compras"
}

// CompraItem guarda un snapshot del nombre y precio del producto al comprar.
type CompraItem struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	CompraID       uint   `gorm:"not null;index" json:"compra_id"`
	ProductoID     uint   `gorm:"not null" json:"producto_id"`
	Nombre         string `gorm:"type:varchar(160);not null" json:"nombre"`
	PrecioUnitario int    `gorm:"not null" json:"precio_unitario"`
	Cantidad       int    `gorm:"not null" json:"cantidad"`
}

func (CompraItem) TableName() string {
	return "compra_items"
}
