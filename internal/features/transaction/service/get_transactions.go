package feature_transaction_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *TransactionService) GetTransactions(
	ctx context.Context,
	filters core_domain.FiltersTransaction,
) ([]core_domain.Transaction, error) {
	if err := filters.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter parameters: %w: %w", err, core_errors.ErrorNotFoud)
	}

	list, err := s.transactionRepository.GetTransactions(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	return list, nil
}
