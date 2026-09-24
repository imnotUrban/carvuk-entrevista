package boilerplate

import (
	"gorm.io/gorm"

	repo "backend/repository/boilerplate"
	service "backend/service/boilerplate"
)

// NewModule is the factory that wires the repository, service and handler for
// the boilerplate feature and returns a ready-to-use Handler.
func NewModule(db *gorm.DB) *Handler {
	return NewHandler(service.NewService(repo.NewRepository(db)))
}
