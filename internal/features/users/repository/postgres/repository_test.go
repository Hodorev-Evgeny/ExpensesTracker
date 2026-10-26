package features_users_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	core_pgx_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
	"github.com/jackc/pgx/v5"
)

type TestRow struct {
	pgx.Row
}

func (r *TestRow) Scan(dest ...interface{}) error {
	return core_repository_pool.ErrNoRows
}

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

type TestData struct {
	name           string
	pool           core_repository_pool.Pool
	function       func(ctx context.Context, repo *UserRepository) (any, error)
	expectedError  error
	expectedResult core_domain.User
}

func TestExtraditionUser(t *testing.T) {
	testPool := CreateTestPool()
	repo := NewUserRepository(testPool)
	testData := []TestData{
		{
			name: "Test get user",
			pool: testPool,
			function: func(ctx context.Context, repo *UserRepository) (any, error) {
				return repo.ExtraditionUser(ctx, -1)
			},
			expectedError:  core_errors.ErrorNotFoud,
			expectedResult: core_domain.User{},
		},
	}

	for _, test := range testData {
		ans, err := test.function(context.Background(), repo)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Test %s expected result %v, but got %v", test.name, test.expectedError, err)
		}
		if ans != test.expectedResult {
			t.Errorf("Test %s expected result %v, but got %v", test.name, test.expectedResult, ans)
		}
	}
}
