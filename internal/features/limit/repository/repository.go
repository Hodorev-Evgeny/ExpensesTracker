package feature_repository_limit

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type LimitRepository struct {
	pool core_repository_pool.Pool
}

func NewLimitRepository(
	pool core_repository_pool.Pool,
) *LimitRepository {
	return &LimitRepository{
		pool: pool,
	}
}
