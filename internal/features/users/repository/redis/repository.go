package feature_user_redis

import (
	core_repository_redis_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/redis/pool"
)

type RepositoryRedis struct {
	rbd core_repository_redis_pool.RedisPool
}

func NewRepositoryRedis(client core_repository_redis_pool.RedisPool) *RepositoryRedis {
	return &RepositoryRedis{
		rbd: client,
	}
}
