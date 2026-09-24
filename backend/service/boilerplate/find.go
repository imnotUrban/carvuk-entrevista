package boilerplate

import (
	"errors"

	"gorm.io/gorm"

	entity "backend/entity/boilerplate"
	"backend/pkg/apperr"
)

func (s *service) GetByID(id uint) (*entity.Boilerplate, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}
	return item, nil
}
