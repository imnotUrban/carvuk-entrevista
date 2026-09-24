package producto

import (
	entity "backend/entity/producto"
	"backend/pkg/dbutil"
)

func (r *repository) List(filter Filter) ([]entity.Producto, int64, error) {
	var items []entity.Producto
	var total int64

	query := r.db.Model(&entity.Producto{})

	if filter.Nombre != "" {
		query = query.Where(dbutil.CaseInsensitiveLike(r.db, "nombre"), "%"+filter.Nombre+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order(filter.SortBy + " " + filter.Order)

	if filter.Limit > 0 {
		offset := 0
		if filter.Page > 1 {
			offset = (filter.Page - 1) * filter.Limit
		}
		query = query.Limit(filter.Limit).Offset(offset)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
