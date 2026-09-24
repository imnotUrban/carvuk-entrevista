// Package dbutil holds small GORM query helpers shared by every repository.
package dbutil

import "gorm.io/gorm"

// CaseInsensitiveLike returns a WHERE clause fragment for a case-insensitive
// partial match on the given column. Postgres supports ILIKE natively;
// other dialects (used in tests, e.g. SQLite) fall back to LOWER(...) LIKE LOWER(?).
func CaseInsensitiveLike(db *gorm.DB, column string) string {
	if db.Dialector.Name() == "postgres" {
		return column + " ILIKE ?"
	}
	return "LOWER(" + column + ") LIKE LOWER(?)"
}
