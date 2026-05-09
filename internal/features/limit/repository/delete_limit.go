package feature_repository_limit

import (
	"context"
	"fmt"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *LimitRepository) DeleteLimit(
	ctx context.Context,
	limitID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		DELETE 
		FROM trackerapp.spending_limit
		WHERE id = $1;
		`

	tag, err := r.pool.Exec(ctx, query, limitID)
	if err != nil {
		return fmt.Errorf("delete limit: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete limit: limit not found: %w", core_errors.ErrorNotFoud)
	}

	return nil
}
