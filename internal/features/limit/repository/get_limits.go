package feature_repository_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *LimitRepository) GetLimits(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]core_domain.Limit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, duration, amount_limit
		FROM trackerapp.spending_limit
		LIMIT $1
		OFFSET $2;`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error scan rows in database: %w", err)
	}
	defer rows.Close()

	limits := make([]core_domain.Limit, 0)
	for rows.Next() {
		l := core_domain.Limit{}
		err := rows.Scan(
			&l.ID,
			&l.Duration,
			&l.AmountLimit,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan row in database: %w", err)
		}
		limits = append(limits, l)
	}

	return limits, nil
}
