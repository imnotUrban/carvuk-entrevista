package boleta

import (
	"gorm.io/gorm"

	boletaentity "backend/entity/boleta"
	compraentity "backend/entity/compra"
	usuarioentity "backend/entity/usuario"
)

// Filter controla la paginación y el orden (por created_at) del listado.
type Filter struct {
	Page  int
	Limit int
	Order string
}

type Repository interface {
	List(filter Filter) ([]boletaentity.Boleta, int64, error)
	GetByID(id uint) (*boletaentity.Boleta, error)
	// Create guarda compra + ítems + boleta y descuenta el stock en una sola transacción.
	Create(compra *compraentity.Compra, boleta *boletaentity.Boleta) error
	FirstUsuario() (*usuarioentity.Usuario, error)
	CountUsuarios() (int64, error)
	CreateUsuario(u *usuarioentity.Usuario) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
