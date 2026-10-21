package feature_service_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *LimitService) DeleteLimit(
	ctx context.Context,
	limitID int,
	userID int,
) error {
	if limitID == core_domain.UnincelizedID {
		return core_errors.ErrorBadRequest
	}

	if err := s.limitRepo.DeleteLimit(ctx, limitID, userID); err != nil {
		return fmt.Errorf("delete limit: %w", err)
	}

	return nil
}
