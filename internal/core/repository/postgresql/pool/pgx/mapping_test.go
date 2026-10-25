package core_pgx_pool

import (
	"errors"
	"testing"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type TestData struct {
	Name        string
	InputCode   *pgconn.PgError
	ExpectedErr error
}

func TestMappingError(t *testing.T) {

	testData := []TestData{
		{
			Name:        "Test mapping Violation Key error",
			InputCode:   &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation},
			ExpectedErr: core_repository_pool.ErrViolationKey,
		},
		{
			Name:        "Test mapping Error Conflict",
			InputCode:   &pgconn.PgError{Code: pgerrcode.UniqueViolation},
			ExpectedErr: core_errors.ErrorConflict,
		},
	}

	for _, test := range testData {
		err := MappingError(test.InputCode)
		if !errors.Is(err, test.ExpectedErr) {
			t.Errorf("MappingError(%v): expected %v, actual %v", test.InputCode, test.ExpectedErr, err)
		}
	}
}
