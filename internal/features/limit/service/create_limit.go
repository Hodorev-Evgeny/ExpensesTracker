package feature_service_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *LimitService) CreateLimit(
	ctx context.Context,
	limit core_domain.Limit,
) (core_domain.Limit, error) {
	if err := limit.Validate(); err != nil {
		return core_domain.Limit{}, fmt.Errorf("message: %w;error: %w", err, core_errors.ErrorValidation)
	}

	limitDomain, err := s.limitRepo.CreateLimit(ctx, limit)
	if err != nil {
		return core_domain.Limit{}, fmt.Errorf("error creating limit in service level: %w", err)
	}
	return limitDomain, nil
}
