package feature_repositor_category

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *CategoryRepository) GetCategoryByID(
	ctx context.Context,
	id int,
) (core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	quey := `
			SELECT id, title, user_id
			FROM trackerapp.categories 
			WHERE id = $1;`

	row := r.pool.QueryRow(ctx, quey, id)

	var category core_domain.Category
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.User_ID,
	)

	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return core_domain.Category{}, fmt.Errorf("category not found: %w", core_errors.ErrorNotFoud)
		}
		return core_domain.Category{}, fmt.Errorf("get category by id: %w", err)
	}

	return category, nil
}
