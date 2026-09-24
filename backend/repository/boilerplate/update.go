package boilerplate

import entity "backend/entity/boilerplate"

func (r *repository) Update(item *entity.Boilerplate) error {
	return r.db.Save(item).Error
}
