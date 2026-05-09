package feature_repositor_category

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *CategoryRepository) GetAllCategories(
	ctx context.Context,
) ([]core_domain.Category, error) {
	query := `
		SELECT id, title, user_id
		FROM trackerapp.categories;`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Error get all categories: %w", err)
	}

	list := make([]core_domain.Category, 0)
	for rows.Next() {
		var category core_domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.User_ID,
		)

		if err != nil {
			return list, fmt.Errorf("Error get list category: %w", err)
		}

		list = append(list, category)
	}

	return list, nil
}
