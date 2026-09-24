package boleta

import boletaentity "backend/entity/boleta"

func (r *repository) GetByID(id uint) (*boletaentity.Boleta, error) {
	var b boletaentity.Boleta
	if err := r.db.Preload("Compra.Usuario").Preload("Compra.Items").First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}
