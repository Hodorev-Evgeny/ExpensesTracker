package feature_repository_health

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *HealthRepository) CheckHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.pool.Ping(ctx); err != nil {
		return fmt.Errorf("postgres health check failed: %w", core_errors.ErrorHealthCheck)
	}

	_, errRedisPing := r.rbd.Ping(ctx).Result()
	if errRedisPing != nil {
		return fmt.Errorf("redis health check failed: %w", core_errors.ErrorHealthCheck)
	}

	return nil
}
