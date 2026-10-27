package core_goredis_pool

import (
	"context"
	"time"

	core_repository_redis_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	ttl     time.Duration
	timeout time.Duration
}

func newRedisClient(config RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ADDR,
		Password: config.Password,
		DB:       config.Database,
	})

	return &RedisClient{
		rdb,
		config.TTL,
		config.TimeOut,
	}, nil
}

func CreateRedisClientMust(config RedisConfig) *RedisClient {
	client, err := newRedisClient(config)
	if err != nil {
		panic(err)
	}

	return client
}

func (r *RedisClient) Ping(ctx context.Context) core_repository_redis_pool.CustomStatusCmd {
	ans := r.Client.Ping(ctx)
	return &CustomStatusCmd{ans}
}

func (r *RedisClient) Get(ctx context.Context, key string) core_repository_redis_pool.CustomStringCmd {
	ans := r.Client.Get(ctx, key)
	return &CustomStringCmd{ans}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) core_repository_redis_pool.CustomStatusCmd {
	ans := r.Client.Set(ctx, key, value, expiration)
	return &CustomStatusCmd{ans}
}

func (r *RedisClient) Del(ctx context.Context, key string) core_repository_redis_pool.CustomIntCmd {
	ans := r.Client.Del(ctx, key)
	return &CustomIntCmd{ans}
}

func (r *RedisClient) GetTTL() time.Duration {
	return r.ttl
}

func (r *RedisClient) GetTimeout() time.Duration {
	return r.timeout
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
