package dbutil

import "gorm.io/gorm"

func CaseInsensitiveLike(db *gorm.DB, column string) string {
	if db.Dialector.Name() == "postgres" {
		return column + " ILIKE ?"
	}
	return "LOWER(" + column + ") LIKE LOWER(?)"
}
