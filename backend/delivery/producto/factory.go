package producto

import (
	"gorm.io/gorm"

	repo "backend/repository/producto"
	service "backend/service/producto"
)

func NewModule(db *gorm.DB) *Handler {
	return NewHandler(service.NewService(repo.NewRepository(db)))
}
