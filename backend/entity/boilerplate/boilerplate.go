package boilerplate

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusDraft    = "draft"
	StatusActive   = "active"
	StatusArchived = "archived"
)

type Boilerplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(160);not null;index:idx_boilerplate_name" json:"name"`
	Code      string         `gorm:"type:varchar(64);not null;index:idx_boilerplate_code" json:"code"`
	Status    string         `gorm:"type:varchar(20);not null;default:'draft';index:idx_boilerplate_status" json:"status"`
	Quantity  int            `gorm:"not null;default:0" json:"quantity"`
	Amount    float64        `gorm:"type:numeric(12,2);not null;default:0" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Boilerplate) TableName() string {
	return "boilerplate"
}
