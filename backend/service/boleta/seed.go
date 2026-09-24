package boleta

import usuarioentity "backend/entity/usuario"

// SeedUsuarioIfEmpty inserta el usuario demo si la tabla está vacía.
func (s *service) SeedUsuarioIfEmpty() error {
	count, err := s.repo.CountUsuarios()
	if err != nil || count > 0 {
		return err
	}
	return s.repo.CreateUsuario(&usuarioentity.Usuario{Nombre: seedUserName, Correo: seedUserCorreo})
}
