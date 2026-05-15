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

func (cmd *CustomStringCmd) Result() (string, error) {
	return cmd.StringCmd.Result()
}

func (status *CustomStatusCmd) Err() error {
	return status.StatusCmd.Err()
}
