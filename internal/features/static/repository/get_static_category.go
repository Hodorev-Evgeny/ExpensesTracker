package feature_repository_static

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *StaticRepository) GetStaticCategories(
	ctx context.Context,
	filters core_domain.FiltersStatic,
) ([]core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, title, user_id
		FROM trackerapp.categories`

	args := []any{}
	conditions := []string{}

	if filters.From != nil {
		conditions = append(conditions, fmt.Sprintf("date>=$%d", len(args)+1))
		args = append(args, *filters.From)
	}
	if filters.To != nil {
		conditions = append(conditions, fmt.Sprintf("date<$%d", len(args)+1))
		args = append(args, *filters.To)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY id ASC"

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
