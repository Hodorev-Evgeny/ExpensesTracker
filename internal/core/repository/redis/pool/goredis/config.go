package core_goredis_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	ADDR     string        `envconfig:"ADDR" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database int           `envconfig:"DATABASE" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"5s"`
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
