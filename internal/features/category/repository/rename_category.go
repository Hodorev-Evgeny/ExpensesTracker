package feature_repositor_category

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *CategoryRepository) RenameCategory(
	ctx context.Context,
	id int,
	category core_domain.Category,
) (core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		UPDATE trackerapp.categories 
		SET title=$1
		WHERE id=$2
		RETURNING id, title, user_id;`

	row := r.pool.QueryRow(ctx, query, category.Name, id)

	var categoryUpdated core_domain.Category
	err := row.Scan(
		&categoryUpdated.ID,
		&categoryUpdated.Name,
		&categoryUpdated.User_ID,
	)
	if err != nil {
		return core_domain.Category{}, err
	}

	return categoryUpdated, nil
}
