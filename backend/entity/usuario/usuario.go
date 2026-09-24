package usuario

import "time"

type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"type:varchar(160);not null" json:"nombre"`
	Correo    string    `gorm:"type:varchar(160);not null" json:"correo"`
	CreatedAt time.Time `json:"created_at"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
