package producto

import entity "backend/entity/producto"

func (r *repository) SoftDelete(id uint) error {
	return r.db.Delete(&entity.Producto{}, id).Error
}
