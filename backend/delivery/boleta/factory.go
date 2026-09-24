package boleta

import (
	"gorm.io/gorm"

	repo "backend/repository/boleta"
	productorepo "backend/repository/producto"
	service "backend/service/boleta"
)

func NewModule(db *gorm.DB) *Handler {
	return NewHandler(service.NewService(repo.NewRepository(db), productorepo.NewRepository(db)))
}
