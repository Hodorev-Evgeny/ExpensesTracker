package feature_repository_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *LimitRepository) CreateLimit(
	ctx context.Context,
	limit core_domain.Limit,
) (core_domain.Limit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO trackerapp.spending_limit(duration, amount_limit)
		VALUES ($1, $2)
		RETURNING id, duration, amount_limit;`

	row := r.pool.QueryRow(ctx, query, limit.Duration, limit.AmountLimit)

	var limitDomain core_domain.Limit
	err := row.Scan(
		&limitDomain.ID,
		&limitDomain.Duration,
		&limitDomain.AmountLimit,
	)
	if err != nil {
		return limitDomain, fmt.Errorf("error inserting limit: %w", err)
	}

	return limitDomain, nil
}
