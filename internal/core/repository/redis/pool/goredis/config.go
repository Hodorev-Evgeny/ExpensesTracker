package core_goredis_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	ADDR     string        `envconfig:"ADDR"`
	Password string        `envconfig:"PASSWORD"`
	Database int           `envconfig:"DATABASE"`
	TTL      time.Duration `envconfig:"TTL"`
}

func getRedisConfig() (RedisConfig, error) {
	redisConfig := RedisConfig{}
	if err := envconfig.Process("REDIS", &redisConfig); err != nil {
		return RedisConfig{}, fmt.Errorf("failed to process redis env var: %w", err)
	}

	return redisConfig, nil
}

func MustGetRedisConfig() RedisConfig {
	redisConfig, err := getRedisConfig()
	if err != nil {
		panic(err)
	}
	return redisConfig
}
