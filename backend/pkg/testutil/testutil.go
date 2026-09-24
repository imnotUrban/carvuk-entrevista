package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	boilerplateentity "backend/entity/boilerplate"
	boletaentity "backend/entity/boleta"
	compraentity "backend/entity/compra"
	productoentity "backend/entity/producto"
	usuarioentity "backend/entity/usuario"
)

func SetupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&boilerplateentity.Boilerplate{}, &productoentity.Producto{},
		&usuarioentity.Usuario{}, &compraentity.Compra{}, &compraentity.CompraItem{}, &boletaentity.Boleta{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}
