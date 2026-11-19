package transactionservice

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	corelogger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"

	"go.uber.org/zap"
)

func (s *TransactionService) DeleteTransaction(
	ctx context.Context,
	id int,
	userID int,
) error {
	if id == core_domain.UninitializedID {
		return core_errors.ErrBadRequest
	}
	log := corelogger.FromContext(ctx)

	transaction, err := s.GetTransaction(ctx, id, userID)
	if err != nil {
		log.Error("error getting transaction for budget:", zap.Error(err))
	}

	if err := s.transactionRepository.DeleteTransaction(ctx, id, userID); err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	if transaction.Type == "Expenditure" {
		ans, err := s.client.RevokeSpending(ctx, core_domain.DeletedTransaction{
			ID:     transaction.ID,
			UserID: userID,
		})
		if err != nil {
			log.Error("error revoking transaction", zap.Error(err))
		} else {
			log.Debug("spending transaction revoked successfully", zap.Any("ans", ans))
		}

	}

	return nil
}
