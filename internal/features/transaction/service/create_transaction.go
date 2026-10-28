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
		return core_domain.Transaction{}, fmt.Errorf("invalid transaction comments: %w: %w",
			err,
			core_errors.ErrorBadRequest,
		)
	}

	exist, err := s.transactionRepository.CheckCategoryOwner(transaction.CategoryID, transaction.UserID)
	if err != nil {
		return core_domain.Transaction{}, fmt.Errorf("failed to check transaction owner: %w", err)
	}
	if !exist {
		return core_domain.Transaction{}, fmt.Errorf("this category dont exists for this user: %w", core_errors.ErrorNotFoud)
	}

	transactionDomain, err := s.transactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return core_domain.Transaction{}, fmt.Errorf("error creating transaction: %w", err)
	}

	return transactionDomain, nil
}
