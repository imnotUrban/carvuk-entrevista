// Package apperr defines the domain-level errors shared across services and
// mapped to HTTP responses by the delivery layer.
package apperr

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrValidation    = errors.New("validation error")
	ErrDuplicateCode = errors.New("a boilerplate with this code already exists")
)
