package categoryrepo

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)


func (r *CategoryRepository) DeleteCategory(
	ctx context.Context,
	id int,
	userID int,
) ([]core_domain.DeletedTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx, core_repository_pool.TxOptions{
		IsolationLevel: core_repository_pool.ReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	queryBlock := `
		SELECT id
		FROM trackerapp.categories
		WHERE id = $1 AND user_id = $2
		FOR UPDATE;
	`

	var categoryID int
	err = tx.QueryRow(ctx, queryBlock, id, userID).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %w", core_errors.ErrNotFound)
		}

		return nil, fmt.Errorf("failed to lock category: %w", err)
	}

	queryDeleteTransactions := `
		DELETE FROM trackerapp.transactions
		WHERE category_id = $1 AND user_id = $2
		RETURNING id, user_id;
	`

	rows, err := tx.Query(ctx, queryDeleteTransactions, categoryID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete transactions: %w", err)
	}

	var deleted []core_domain.DeletedTransaction

	for rows.Next() {
		var transaction core_domain.DeletedTransaction

		if err := rows.Scan(&transaction.ID, &transaction.UserID); err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan deleted transaction: %w", err)
		}

		deleted = append(deleted, transaction)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read deleted transactions: %w", err)
	}

	queryDeleteCategory := `
		DELETE FROM trackerapp.categories
		WHERE id = $1 AND user_id = $2;
	`

	tag, err := tx.Exec(ctx, queryDeleteCategory, categoryID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return nil, fmt.Errorf("category not found: %w", core_errors.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return deleted, nil
}