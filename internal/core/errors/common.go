package core_errors

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrValidation   = errors.New("validation error")
	ErrHealthCheck  = errors.New("health check failed")
	ErrConflict     = errors.New("already exists")
)
