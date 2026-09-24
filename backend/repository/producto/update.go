package producto

import entity "backend/entity/producto"

func (r *repository) Update(item *entity.Producto) error {
	return r.db.Save(item).Error
}
