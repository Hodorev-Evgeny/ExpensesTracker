package core_goredis_pool

import (
	"github.com/redis/go-redis/v9"
)

type StringCmd struct {
	*redis.StringCmd
}

type StatusCmd struct {
	*redis.StatusCmd
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.StringCmd.Result()
}

func (status *StatusCmd) Err() error {
	return status.StatusCmd.Err()
}
