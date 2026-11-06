package core_errors

import (
	"errors"
)

var (
	ErrorNotFoud      = errors.New("not foud")
	ErrorBadRequest   = errors.New("bad request")
	ErrorUnauthorized = errors.New("unauthorized")
	ErrorValidation   = errors.New("validation error")
	ErrorHealthCheck  = errors.New("health check failed")
	ErrorConflict     = errors.New("already exists")
)
