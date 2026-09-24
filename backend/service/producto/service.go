package producto

import (
	"fmt"

	entity "backend/entity/producto"
	"backend/pkg/apperr"
	repo "backend/repository/producto"
)

var sortWhitelist = map[string]bool{
	"id":         true,
	"nombre":     true,
	"precio":     true,
	"stock":      true,
	"created_at": true,
}

type Service interface {
	Create(input CreateInput) (*entity.Producto, error)
	GetByID(id uint) (*entity.Producto, error)
	List(input ListInput) ([]entity.Producto, int64, error)
	Update(id uint, input UpdateInput) (*entity.Producto, error)
	Delete(id uint) error
	SeedIfEmpty() error
}

type service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) Service {
	return &service{repo: repo}
}

func validatePrecio(precio int) error {
	if precio <= 0 {
		return fmt.Errorf("%w: precio must be greater than 0", apperr.ErrValidation)
	}
	return nil
}

func validateStock(stock int) error {
	if stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", apperr.ErrValidation)
	}
	return nil
}
