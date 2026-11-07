package transactionservice

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *TransactionService) DeleteTransaction(
	ctx context.Context,
	id int,
	userID int,
) error {
	if id == core_domain.UninitializedID {
		return core_errors.ErrBadRequest
	}

	if err := s.transactionRepository.DeleteTransaction(ctx, id, userID); err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	return nil
}
