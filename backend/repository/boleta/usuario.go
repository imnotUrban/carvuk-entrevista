package boleta

import usuarioentity "backend/entity/usuario"

func (r *repository) FirstUsuario() (*usuarioentity.Usuario, error) {
	var u usuarioentity.Usuario
	if err := r.db.Order("id ASC").First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) CountUsuarios() (int64, error) {
	var count int64
	if err := r.db.Model(&usuarioentity.Usuario{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) CreateUsuario(u *usuarioentity.Usuario) error {
	return r.db.Create(u).Error
}
