package feature_repository_transaction

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
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
		return core_domain.Transaction{}, fmt.Errorf("Error scanning transaction: %w", err)
	}

	return transactionDomain, nil
}
