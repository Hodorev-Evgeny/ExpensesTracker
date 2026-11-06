package core_repository_pool

import "errors"

var (
	ErrNoRows       = errors.New("no rows")
	ErrViolationKey = errors.New("violation key error")
)
