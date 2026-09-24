// Package testutil provides a shared in-memory SQLite database for the
// service-layer test suites (service/boilerplate), so they run fast and
// without a real PostgreSQL dependency.
package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	boilerplateentity "backend/entity/boilerplate"
)

// SetupDB spins up an isolated in-memory SQLite database with the same
// schema as production (via AutoMigrate).
func SetupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&boilerplateentity.Boilerplate{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}
