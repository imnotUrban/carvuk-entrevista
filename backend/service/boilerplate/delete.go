package boilerplate

func (s *service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}
