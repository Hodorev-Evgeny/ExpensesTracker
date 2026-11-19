package categoryservice

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	corelogger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"go.uber.org/zap"
)

func (s *CategoryService) DeleteCategory(
	ctx context.Context,
	id int,
	userId int,
) error {
	if id == core_domain.UninitializedID {
		return fmt.Errorf("id category is invalid")
	}

	listTransaction, status := s.repository.DeleteCategory(ctx, id, userId)

	log := corelogger.FromContext(ctx)
	for _, transaction := range listTransaction {
		ans, err := s.client.RevokeSpending(ctx, transaction)
		if err != nil {
			log.Error("Error revoking transaction", zap.Error(err))
			continue
		}
		log.Debug("correct revoke transaction", zap.Any("ans", ans))
	}

	return status
}
