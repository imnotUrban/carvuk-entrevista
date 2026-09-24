package boilerplate

import (
	entity "backend/entity/boilerplate"
	"backend/pkg/dbutil"
)

func (r *repository) List(filter Filter) ([]entity.Boilerplate, int64, error) {
	var items []entity.Boilerplate
	var total int64

	query := r.db.Model(&entity.Boilerplate{})

	if filter.Name != "" {
		query = query.Where(dbutil.CaseInsensitiveLike(r.db, "name"), "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		query = query.Where(dbutil.CaseInsensitiveLike(r.db, "code"), "%"+filter.Code+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	order := filter.Order
	if order == "" {
		order = "asc"
	}
	query = query.Order(sortBy + " " + order)

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
