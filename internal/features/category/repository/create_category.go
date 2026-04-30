package feature_repositor_category

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *CategoryRepository) CreateCategory(
	ctx context.Context,
	category core_domain.Category,
) (core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO trackerapp.categories (title, user_id)
		VALUES ($1, $2) 
		RETURNING id, title, user_id;`

	row := r.pool.QueryRow(ctx,
		query,
		category.Name,
		category.User_ID,
	)

	err := row.Scan(&category.ID,
		&category.Name,
		&category.User_ID,
	)
	if err != nil {
		return core_domain.Category{}, fmt.Errorf("create category: %w", err)
	}

	return category, nil
}
