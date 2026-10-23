package mock_repository

import (
	"context"
	"time"

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

func (c *Pool) Query(ctx context.Context, sql string, args ...any) (core_repository_pool.Rows, error) {
	return core_pgx_pool.Rows{}, nil
}

func (c *Pool) Close() {
}

func CreateTestPool() *Pool {
	return &Pool{}
}
