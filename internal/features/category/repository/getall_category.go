package feature_repositor_category

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *CategoryRepository) GetAllCategories(
	ctx context.Context,
	userId int,
) ([]core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, title, user_id
		FROM trackerapp.categories
		WHERE user_id = $1;`

	rows, err := r.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("Error get all categories: %w", err)
	}
	defer rows.Close()

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

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}

	return list, nil
}
