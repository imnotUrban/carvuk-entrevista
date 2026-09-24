package producto

import (
	"gorm.io/gorm"

	entity "backend/entity/producto"
)

type Repository interface {
	Create(item *entity.Producto) error
	GetByID(id uint) (*entity.Producto, error)
	List(filter Filter) ([]entity.Producto, int64, error)
	Update(item *entity.Producto) error
	SoftDelete(id uint) error
	CountAll() (int64, error)
}

type Filter struct {
	Nombre string
	Page   int
	Limit  int
	SortBy string
	Order  string
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
