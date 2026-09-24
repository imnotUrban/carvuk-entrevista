package boilerplate

import (
	"gorm.io/gorm"

	repo "backend/repository/boilerplate"
	service "backend/service/boilerplate"
)

func NewModule(db *gorm.DB) *Handler {
	return NewHandler(service.NewService(repo.NewRepository(db)))
}
