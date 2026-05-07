package feature_service_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *LimitService) GetLimit(
	ctx context.Context,
	id int,
) (core_domain.Limit, error) {
	if id == core_domain.UnincelizedID {
		return core_domain.Limit{}, core_errors.ErrorUnauthorized
	}

	limitDomain, err := s.limitRepo.GetLimit(ctx, id)
	if err != nil {
		return core_domain.Limit{}, fmt.Errorf("error getting in service: %w", err)
	}

	return limitDomain, nil
}
