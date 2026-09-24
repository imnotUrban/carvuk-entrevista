package producto

import entity "backend/entity/producto"

func (r *repository) GetByID(id uint) (*entity.Producto, error) {
	var item entity.Producto
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// CountAll cuenta todas las filas, incluidas las soft-deleteadas.
func (r *repository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Unscoped().Model(&entity.Producto{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
