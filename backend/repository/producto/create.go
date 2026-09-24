package producto

import entity "backend/entity/producto"

func (r *repository) Create(item *entity.Producto) error {
	return r.db.Create(item).Error
}
