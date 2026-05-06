package feature_repository_transaction

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *TransactionRepository) PatchTransaction(
	ctx context.Context,
	transaction core_domain.Transaction,
) (core_domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	quary := `
		UPDATE trackerapp.transactions
		SET sum=$1, category_id=$2, type_transaction=$3, comments=$4, time_changes=$5
		WHERE id=$6
		RETURNING id, sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes;`

	row := r.pool.QueryRow(ctx, quary,
		transaction.Sum,
		transaction.CategoryID,
		transaction.Type,
		transaction.Comments,
		transaction.TimeChange,
		transaction.ID,
	)

	var updatedTransaction core_domain.Transaction
	err := row.Scan(
		&updatedTransaction.ID,
		&updatedTransaction.Sum,
		&updatedTransaction.Type,
		&updatedTransaction.Date,
		&updatedTransaction.CategoryID,
		&updatedTransaction.UserID,
		&updatedTransaction.Comments,
		&updatedTransaction.TimeCreated,
		&updatedTransaction.TimeChange,
	)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrViolationKey) {
			return core_domain.Transaction{}, fmt.Errorf("error key violation: %w", core_errors.ErrorValidation)
		}
		return updatedTransaction, fmt.Errorf("error updating transaction: %w", err)
	}

	return updatedTransaction, nil
}
