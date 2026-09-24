package boilerplate

import entity "backend/entity/boilerplate"

func (r *repository) GetByID(id uint) (*entity.Boilerplate, error) {
	var item entity.Boilerplate
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) ExistsByCode(code string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&entity.Boilerplate{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
