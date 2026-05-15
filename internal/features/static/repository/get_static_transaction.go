package feature_repository_static

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *StaticRepository) GetStaticTransactions(
	ctx context.Context,
	filters core_domain.FiltersStatic,
) ([]core_domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes
		FROM trackerapp.transactions
		`

	args := []any{}
	conditions := []string{}

	conditions = append(conditions, fmt.Sprintf("user_id=$%d", len(args)+1))
	args = append(args, filters.UserID)

	if filters.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id=$%d", len(args)+1))
		args = append(args, *filters.CategoryID)
	}
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

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Error get all categories: %w", err)
	}

	var list []core_domain.Transaction
	for rows.Next() {
		var transaction core_domain.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.Sum,
			&transaction.Type,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.UserID,
			&transaction.Comments,
			&transaction.TimeCreated,
			&transaction.TimeChange,
		)

		if err != nil {
			return nil, fmt.Errorf("transaction repository GetTransactions: %w", err)
		}
		list = append(list, transaction)
	}

	return list, nil
}
