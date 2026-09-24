package boilerplate

// Delete only soft-deletes: it never removes the row physically (see
// repository/boilerplate/delete.go).
func (s *service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}
