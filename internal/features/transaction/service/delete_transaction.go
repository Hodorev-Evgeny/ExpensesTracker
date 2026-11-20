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

	if err := s.transactionRepository.DeleteTransaction(ctx, id, userID); err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	ans, err := s.client.RevokeSpending(ctx, core_domain.DeletedTransaction{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		log.Error("error revoking transaction", zap.Error(err))
	} else {
		log.Debug("spending transaction revoked successfully", zap.Any("ans", ans))
	}

	return nil
}
