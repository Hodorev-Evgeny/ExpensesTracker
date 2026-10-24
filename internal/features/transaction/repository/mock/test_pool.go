package transaction_repository_mock

import (
	"context"
	"strings"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	core_pgx_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
)

type Pool struct {
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) core_repository_pool.Row {
	return &TestRow{}
}

func (p *Pool) GetTimeout() time.Duration {
	return time.Second
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (core_repository_pool.CommandTag, error) {
	return core_pgx_pool.CommandTag{}, nil
}

func (p *Pool) Ping(ctx context.Context) error {
	return nil
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (core_repository_pool.Rows, error) {
	if strings.Contains(sql, "limit") {
		if !strings.Contains(sql, "limit=$") {
			return core_pgx_pool.Rows{}, core_errors.ErrorNotPlaceholder
		}
	}

	if strings.Contains(sql, "offset") {
		if !strings.Contains(sql, "offset=$") {
			return core_pgx_pool.Rows{}, core_errors.ErrorNotPlaceholder
		}
	}

	return core_pgx_pool.Rows{}, core_errors.ErrorMock
}

func (p *Pool) Close() {

}

func CreateTestPool() *Pool {
	return &Pool{}
}
