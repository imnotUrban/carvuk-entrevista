package producto

import entity "backend/entity/producto"

var seedProductos = []entity.Producto{
	{Nombre: "Leche", Precio: 1100, Stock: 50},
	{Nombre: "Pan", Precio: 1800, Stock: 100},
	{Nombre: "Aceite", Precio: 2000, Stock: 30},
}

// SeedIfEmpty inserta los productos de ejemplo solo si la tabla no tiene filas.
func (s *service) SeedIfEmpty() error {
	count, err := s.repo.CountAll()
	if err != nil || count > 0 {
		return err
	}
	for i := range seedProductos {
		item := seedProductos[i]
		if err := s.repo.Create(&item); err != nil {
			return err
		}
	}
	return nil
}
