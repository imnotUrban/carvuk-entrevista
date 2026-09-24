package producto

import (
	"time"

	"gorm.io/gorm"
)

type Producto struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nombre    string         `gorm:"type:varchar(160);not null;index:idx_productos_nombre" json:"nombre"`
	Precio    int            `gorm:"not null" json:"precio"`
	Stock     int            `gorm:"not null;default:0" json:"stock"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Producto) TableName() string {
	return "productos"
}
