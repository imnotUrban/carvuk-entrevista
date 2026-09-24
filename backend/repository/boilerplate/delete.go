package boilerplate

import entity "backend/entity/boilerplate"

// SoftDelete never removes the row physically: GORM turns this into
// `UPDATE boilerplate SET deleted_at = now() WHERE id = ?` because Boilerplate
// embeds gorm.DeletedAt.
func (r *repository) SoftDelete(id uint) error {
	return r.db.Delete(&entity.Boilerplate{}, id).Error
}
