package boilerplate

import entity "backend/entity/boilerplate"

func (r *repository) Create(item *entity.Boilerplate) error {
	return r.db.Create(item).Error
}
