package feature_repository_session

import core_repository_redis_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/redis/pool"

type SessionRepository struct {
	rbd core_repository_redis_pool.RedisPool
}

func NewSessionRepository(
	rbd core_repository_redis_pool.RedisPool,
) *SessionRepository {
	return &SessionRepository{
		rbd: rbd,
	}
}
