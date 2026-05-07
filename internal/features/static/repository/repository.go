package feature_repository_static

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type StaticRepository struct {
	pool core_repository_pool.Pool
}

func NewStaticRepository(
	pool core_repository_pool.Pool,
) *StaticRepository {
	return &StaticRepository{
		pool: pool,
	}
}
