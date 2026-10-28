package feature_repository_transaction

import (
	"context"
	"fmt"
)

func (r *TransactionRepository) CheckCategoryOwner(
	ctx context.Context,
	categoryID int,
	userID int,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1 FROM trackerapp.categories WHERE id = $1 AND user_id = $2);`

	var exists bool
	err := r.pool.QueryRow(ctx, query, categoryID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check category owner: %w", err)
	}

	return exists, nil
}
