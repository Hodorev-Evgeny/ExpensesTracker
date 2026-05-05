package feature_transaction_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *TransactionService) PatchTransaction(
	ctx context.Context,
	id int,
	patch core_domain.PatchTransaction,
) (core_domain.Transaction, error) {
	transaction, err := s.GetTransaction(ctx, id)
	if err != nil {
		return core_domain.Transaction{}, fmt.Errorf("error getting transaction: %w", err)
	}

	if err := transaction.ApplyPatch(patch); err != nil {
		return core_domain.Transaction{}, fmt.Errorf("error patching transaction: %w", err)
	}

	updateTransaction, err := s.transactionRepository.PatchTransaction(ctx, transaction)
	if err != nil {
		return core_domain.Transaction{}, fmt.Errorf("error updating transaction: %w", err)
	}

	return updateTransaction, nil
}
