package feature_repository_transaction

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *TransactionRepository) CreateTransaction(
	ctx context.Context,
	transaction core_domain.Transaction,
) (core_domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO trackerapp.transactions (sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes`

	row := r.pool.QueryRow(ctx, query,
		transaction.Sum,
		transaction.Type,
		transaction.Date,
		transaction.CategoryID,
		transaction.UserID,
		transaction.Comments,
		transaction.TimeCreated,
		transaction.TimeChange,
	)

	var transactionDomain core_domain.Transaction
	err := row.Scan(&transactionDomain.ID,
		&transactionDomain.Sum,
		&transactionDomain.Type,
		&transactionDomain.Date,
		&transactionDomain.CategoryID,
		&transactionDomain.UserID,
		&transactionDomain.Comments,
		&transactionDomain.TimeCreated,
		&transactionDomain.TimeChange,
	)

	if err != nil {
		if errors.Is(err, core_repository_pool.ErrViolationKey) {
			return core_domain.Transaction{},
				fmt.Errorf("cant create transaction: category not found: %w", core_errors.ErrorBadRequest)
		}
		return core_domain.Transaction{}, fmt.Errorf("Error scanning transaction: %w", err)
	}

	return transactionDomain, nil
}
