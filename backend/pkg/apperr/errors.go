package apperr

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrValidation    = errors.New("validation error")
	ErrDuplicateCode = errors.New("a boilerplate with this code already exists")
)
