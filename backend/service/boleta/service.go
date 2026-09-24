package boleta

import (
	"time"

	usuarioentity "backend/entity/usuario"
	repo "backend/repository/boleta"
	productorepo "backend/repository/producto"
)

const (
	maxLines       = 100
	maxCantidad    = 99
	seedUserName   = "Cliente Demo"
	seedUserCorreo = "cliente.demo@example.com"
)

type Service interface {
	Create(input CreateInput) (*Detail, error)
	List(input ListInput) ([]ListItem, int64, error)
	GetByID(id uint) (*Detail, error)
	SeedUsuarioIfEmpty() error
}

type service struct {
	repo      repo.Repository
	productos productorepo.Repository
}

func NewService(repo repo.Repository, productos productorepo.Repository) Service {
	return &service{repo: repo, productos: productos}
}

// Detail es la boleta emitida junto con su usuario e ítems.
type Detail struct {
	ID                 uint                   `json:"id"`
	CompraID           uint                   `json:"compra_id"`
	ValorBruto         int                    `json:"valor_bruto"`
	Impuesto           int                    `json:"impuesto"`
	ValorNeto          int                    `json:"valor_neto"`
	PorcentajeImpuesto int                    `json:"porcentaje_impuesto"`
	Usuario            *usuarioentity.Usuario `json:"usuario"`
	Items              []DetailItem           `json:"items"`
	CreatedAt          time.Time              `json:"created_at"`
}

type DetailItem struct {
	ProductoID     uint   `json:"producto_id"`
	Nombre         string `json:"nombre"`
	PrecioUnitario int    `json:"precio_unitario"`
	Cantidad       int    `json:"cantidad"`
	Subtotal       int    `json:"subtotal"`
}
