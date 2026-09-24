package boleta

import (
	boletaentity "backend/entity/boleta"
)

func (r *repository) List(filter Filter) ([]boletaentity.Boleta, int64, error) {
	var items []boletaentity.Boleta
	var total int64

	if err := r.db.Model(&boletaentity.Boleta{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// id como desempate para que la paginación sea estable.
	query := r.db.Preload("Compra.Usuario").Preload("Compra.Items").
		Order("created_at " + filter.Order).Order("id " + filter.Order)
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
