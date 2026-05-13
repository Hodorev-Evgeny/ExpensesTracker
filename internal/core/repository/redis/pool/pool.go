package core_repository_redis_pool

import (
	"context"
	"time"
)

type RedisPool interface {
	Ping(ctx context.Context) *StatusCmd
	Get(ctx context.Context, key string) *StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	Close() error
}

type StatusCmd interface {
	Err() error
}

type StringCmd interface {
	Result() (string, error)
}
