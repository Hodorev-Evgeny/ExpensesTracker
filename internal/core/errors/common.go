package core_errors

import (
	"errors"
)

var (
	// ошибки для обычный запросов
	ErrorNotFoud      = errors.New("Not foud")
	ErrorBadRequest   = errors.New("Bad request")
	ErrorUnauthorized = errors.New("Unauthorized")
	ErrorValidation   = errors.New("Validation error")

	// ошибки для тестов
	ErrorNotPlaceholder = errors.New("Not placeholder")
	ErrorMock           = errors.New("Mock error")
)
