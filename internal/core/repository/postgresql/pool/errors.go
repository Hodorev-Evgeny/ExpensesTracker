package core_repository_pool

import "errors"

var ErrNoRows = errors.New("no rows")
var ErrViolationKey = errors.New("violation key error")
var ErrErrConflict = errors.New("violation error")
