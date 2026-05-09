package feature_repository_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *LimitRepository) PatchLimit(
	ctx context.Context,
	limit core_domain.Limit,
) (core_domain.Limit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		UPDATE trackerapp.spending_limit
		SET duration=$1, amount_limit=$2
		WHERE id=$3
		RETURNING id, duration, amount_limit;`

	row := r.pool.QueryRow(ctx, query, limit.Duration, limit.AmountLimit, limit.ID)

	var updatedLimit core_domain.Limit
	if err := row.Scan(&updatedLimit.ID, &updatedLimit.Duration, &updatedLimit.AmountLimit); err != nil {
		return core_domain.Limit{}, fmt.Errorf("failed to update limit: %w", err)
	}

	return updatedLimit, nil
}
