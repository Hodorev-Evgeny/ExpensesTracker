package core_repository_redis_pool

import (
	"context"
	"time"
)

type RedisPool interface {
	Ping(ctx context.Context) CustomStatusCmd
	Get(ctx context.Context, key string) CustomStringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) CustomStatusCmd
	Close() error
}

type CustomStatusCmd interface {
	Err() error
}

type CustomStringCmd interface {
	Result() (string, error)
}
