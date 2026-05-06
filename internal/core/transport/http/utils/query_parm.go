package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func GetIntQueryParm(r *http.Request, key string) (*int, error) {
	value := r.URL.Query().Get(key)

	if value == "" {
		return nil, nil
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value for value %s: %e", key, core_errors.ErrorValidation)
	}

	return &valueInt, nil
}

func GetDateQueryParm(r *http.Request, key string) (*time.Time, error) {
	value := r.URL.Query().Get(key)

	if value == "" {
		return nil, nil
	}

	valueDate, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return nil, fmt.Errorf("invalid value for value %s: %e", key, core_errors.ErrorValidation)
	}

	return &valueDate, nil
}
