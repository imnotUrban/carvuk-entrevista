package boilerplate

import entity "backend/entity/boilerplate"

func (r *repository) SoftDelete(id uint) error {
	return r.db.Delete(&entity.Boilerplate{}, id).Error
}
