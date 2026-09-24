// Package boilerplate is the service layer for boilerplate rows: validation
// and business rules live here, split one operation per file.
package boilerplate

import (
	"fmt"

	entity "backend/entity/boilerplate"
	"backend/pkg/apperr"
	repo "backend/repository/boilerplate"
)

var sortWhitelist = map[string]bool{
	"id":         true,
	"name":       true,
	"amount":     true,
	"quantity":   true,
	"created_at": true,
	"updated_at": true,
}

var statusWhitelist = map[string]bool{
	entity.StatusDraft:    true,
	entity.StatusActive:   true,
	entity.StatusArchived: true,
}

type Service interface {
	Create(input CreateInput) (*entity.Boilerplate, error)
	GetByID(id uint) (*entity.Boilerplate, error)
	List(input ListInput) ([]entity.Boilerplate, int64, error)
	Update(id uint, input UpdateInput) (*entity.Boilerplate, error)
	Delete(id uint) error
}

type service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) Service {
	return &service{repo: repo}
}

func validateStatus(status string) error {
	if !statusWhitelist[status] {
		return fmt.Errorf("%w: status must be one of draft, active, archived", apperr.ErrValidation)
	}
	return nil
}
