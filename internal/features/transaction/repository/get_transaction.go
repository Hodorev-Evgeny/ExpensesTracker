package feature_repository_transaction

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *TransactionRepository) GetTransaction(
	ctx context.Context,
	id int,
) (core_domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes
		FROM trackerapp.transactions
		WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, id)

	var transaction core_domain.Transaction
	err := row.Scan(
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
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return core_domain.Transaction{}, fmt.Errorf("transaction not found: %w", core_errors.ErrorNotFoud)
		}
		return core_domain.Transaction{}, fmt.Errorf("error getting transaction: %w", err)
	}

	return transaction, nil
}
