package feature_repository_transaction

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *TransactionRepository) DeleteTransaction(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		DELETE 
		FROM trackerapp.transactions
		WHERE id = $1;
		`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete transaction: transaction not found: %w", core_errors.ErrorNotFoud)
	}

	return nil
}
