package feature_repositor_category

import (
	"context"
	"fmt"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *CategoryRepository) DeleteCategory(
	ctx context.Context,
	id int,
	userId int,
) error {
	query := `
		DELETE 
		FROM trackerapp.categories 
		WHERE id = $1 AND user_id = $2;`

	tag, err := r.pool.Exec(ctx, query, id, userId)
	if err != nil {
		return fmt.Errorf(`user with id %d: %w`, id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("id category not in database: %d, %w", id, core_errors.ErrorNotFoud)
	}

	return nil
}
