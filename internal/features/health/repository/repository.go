package feature_repository_health

import (
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	core_repository_redis_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/redis/pool"
)

type HealthRepository struct {
	rbd  core_repository_redis_pool.RedisPool
	pool core_repository_pool.Pool
}

func NewHealthRepository(
	rbd core_repository_redis_pool.RedisPool,
	pool core_repository_pool.Pool,
) *HealthRepository {
	return &HealthRepository{
		rbd:  rbd,
		pool: pool,
	}
}
