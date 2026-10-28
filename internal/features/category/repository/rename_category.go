package feature_repositor_category

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *CategoryRepository) RenameCategory(
	ctx context.Context,
	id int,
	userId int,
	category core_domain.Category,
) (core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		UPDATE trackerapp.categories 
		SET title=$1::text, limit_id=$2::integer
		WHERE id=$3::integer
		  AND user_id = $4::integer
		  AND ($2::integer IS NULL OR EXISTS (
			  SELECT 1 
			  FROM trackerapp.spending_limit 
			  WHERE id = $2::integer AND user_id = $4::integer
		))
		RETURNING id, title, user_id, limit_id;`

	row := r.pool.QueryRow(ctx, query, category.Name, category.Limit_id, id, userId)

	var categoryUpdated core_domain.Category
	err := row.Scan(
		&categoryUpdated.ID,
		&categoryUpdated.Name,
		&categoryUpdated.User_ID,
		&categoryUpdated.Limit_id,
	)
	if err != nil {
		return core_domain.Category{}, err
	}

	return categoryUpdated, nil
}
