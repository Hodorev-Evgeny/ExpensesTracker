package feature_service_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *LimitService) GetLimits(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]core_domain.Limit, error) {
	if (limit != nil && *limit < 0) || (offset != nil && *offset < 0) {
		return []core_domain.Limit{}, fmt.Errorf("invalid limit or offset: %w", core_errors.ErrorValidation)
	}

	list, err := s.limitRepo.GetLimits(ctx, limit, offset)
	if err != nil {
		return []core_domain.Limit{}, fmt.Errorf("limit service error: %w", err)
	}

	return list, nil
}
