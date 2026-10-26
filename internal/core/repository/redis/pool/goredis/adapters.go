package core_goredis_pool

import (
	"github.com/redis/go-redis/v9"
)

type CustomStringCmd struct {
	*redis.StringCmd
}

type CustomStatusCmd struct {
	*redis.StatusCmd
}

type CustomIntCmd struct {
	*redis.IntCmd
}

func (cmd *CustomStringCmd) Result() (string, error) {
	return cmd.StringCmd.Result()
}

func (status *CustomStatusCmd) Err() error {
	return status.StatusCmd.Err()
}

func (status *CustomStatusCmd) Result() (string, error) {
	return status.StatusCmd.Result()
}

func (c *CustomIntCmd) Err() error {
	return c.IntCmd.Err()
}
