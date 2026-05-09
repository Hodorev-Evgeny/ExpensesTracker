package core_pgx_pool

import (
	"errors"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Row struct {
	pgx.Row
}

type Rows struct {
	pgx.Rows
}

func (r Row) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		newError := MappingError(err)
		return newError
	}
	return nil
}

type CommandTag struct {
	pgconn.CommandTag
}

func MappingError(err error) error {
	var pgxErr *pgconn.PgError

	if errors.Is(err, pgx.ErrNoRows) {
		return core_repository_pool.ErrNoRows
	}

	if errors.As(err, &pgxErr) {
		if pgxErr.Code == "23503" {
			return core_repository_pool.ErrViolationKey
		}
	}

	return err
}
