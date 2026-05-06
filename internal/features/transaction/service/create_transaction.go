package feature_transaction_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	transaction core_domain.Transaction,
) (core_domain.Transaction, error) {
	if err := transaction.Validate(); err != nil {
		return core_domain.Transaction{}, fmt.Errorf("invalid transaction comments: %d: %d",
			err,
			core_errors.ErrorNotFoud,
		)
	}

	transactionDomain, err := s.transactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return core_domain.Transaction{}, fmt.Errorf("error creating transaction: %w", err)
	}

	return transactionDomain, nil
}
