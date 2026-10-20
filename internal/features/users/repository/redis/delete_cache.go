package feature_user_redis

import (
	"context"
	"fmt"
	"time"
)

func (r *RepositoryRedis) DeleteCache(
	ctx context.Context,
	sessionID string,
) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := "sessionID:" + sessionID
	cmd := r.rbd.Del(ctx, key)
	if cmd.Err() != nil {
		return fmt.Errorf("delete cache: %w", cmd.Err())
	}
	return nil
}
