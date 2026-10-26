package transaction_repository_mock

import (
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgx/v5"
)

type TestRow struct {
	pgx.Row
}

func (r *TestRow) Scan(dest ...interface{}) error {
	return core_repository_pool.ErrNoRows
}

type emptyRows struct {
}

func (r *emptyRows) Close() {
}

func (r *emptyRows) Err() error {
	return nil
}

func (r *emptyRows) Next() bool {
	return false
}
func (r *emptyRows) Scan(dest ...interface{}) error {
	return nil
}
