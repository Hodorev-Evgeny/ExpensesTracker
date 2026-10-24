package transaction_repository_mock

import (
	"context"
	"time"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

type TestPool struct {
	sql  string
	args []any
}

func (p *TestPool) Query(ctx context.Context, sql string, args ...any) (core_repository_pool.Rows, error) {
	p.sql, p.args = sql, args
	return &emptyRows{}, nil
}

func (p *TestPool) GetQueryAndArgs() (string, []any) {
	return p.sql, p.args
}

func (p *TestPool) QueryRow(ctx context.Context, sql string, args ...any) core_repository_pool.Row {
	return nil
}

func (p *TestPool) Exec(ctx context.Context, sql string, arguments ...any) (core_repository_pool.CommandTag, error) {
	return nil, nil
}

func (p *TestPool) GetTimeout() time.Duration {
	return time.Second
}

func (p *TestPool) Close() {}

func CreateTestPool() *TestPool {
	return &TestPool{}
}
