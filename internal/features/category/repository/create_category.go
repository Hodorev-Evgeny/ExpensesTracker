package feature_repositor_category

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *CategoryRepository) CreateCategory(
	ctx context.Context,
	category core_domain.Category,
) (core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO trackerapp.categories (title, user_id, limit_id)
		VALUES ($1, $2, $3) 
		RETURNING id, title, user_id, limit_id;`

	row := r.pool.QueryRow(ctx,
		query,
		category.Name,
		category.User_ID,
		category.Limit_id,
	)

	err := row.Scan(&category.ID,
		&category.Name,
		&category.User_ID,
		&category.Limit_id,
	)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrViolationKey) {
			return core_domain.Category{}, fmt.Errorf("this limit not in database: %w", core_errors.ErrorBadRequest)
		}
		return core_domain.Category{}, fmt.Errorf("create category: %w", err)
	}

	return category, nil
}
