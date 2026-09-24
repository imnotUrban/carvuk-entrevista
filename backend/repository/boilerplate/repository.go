package boilerplate

import (
	"gorm.io/gorm"

	entity "backend/entity/boilerplate"
)

type Repository interface {
	Create(item *entity.Boilerplate) error
	GetByID(id uint) (*entity.Boilerplate, error)
	List(filter Filter) ([]entity.Boilerplate, int64, error)
	Update(item *entity.Boilerplate) error
	SoftDelete(id uint) error
	ExistsByCode(code string, excludeID uint) (bool, error)
}

type Filter struct {
	Name      string
	Code      string
	Status    string
	MinAmount *float64
	MaxAmount *float64
	Page      int
	Limit     int
	SortBy    string
	Order     string
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
