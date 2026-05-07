package feature_repository_limit

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *LimitRepository) GetLimit(
	ctx context.Context,
	id int,
) (core_domain.Limit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, duration, amount_limit
		FROM trackerapp.spending_limit
		WHERE id=$1;`

	row := r.pool.QueryRow(ctx, query, id)

	var limitDomain core_domain.Limit
	err := row.Scan(
		&limitDomain.ID,
		&limitDomain.Duration,
		&limitDomain.AmountLimit,
	)

	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return core_domain.Limit{}, fmt.Errorf("limit not in database: %w", core_errors.ErrorNotFoud)
		}
		return core_domain.Limit{}, fmt.Errorf("error scanning row: %w", err)
	}

	return limitDomain, nil
}
