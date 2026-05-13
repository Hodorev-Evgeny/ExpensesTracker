package core_goredis_pool

import (
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func newRedisClient(config RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ADDR,
		Password: config.Password,
		DB:       config.Database,
	})

	return &RedisClient{
		rdb,
	}, nil
}

func CreateRedisClientMust(config RedisConfig) *RedisClient {
	client, err := newRedisClient(config)
	if err != nil {
		panic(err)
	}

	return client
}
