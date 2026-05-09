package feature_transaction_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *TransactionService) GetTransaction(
	ctx context.Context,
	id int,
) (core_domain.Transaction, error) {
	if id == core_domain.UnincelizedID {
		return core_domain.Transaction{}, fmt.Errorf("incorrect transaction id")
	}

	transaction, err := s.transactionRepository.GetTransaction(ctx, id)
	if err != nil {
		return core_domain.Transaction{}, fmt.Errorf("error getting transaction: %w", err)
	}

	return transaction, nil
}
