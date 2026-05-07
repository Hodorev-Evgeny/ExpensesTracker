package feature_service_limit

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *LimitService) PatchLimit(
	ctx context.Context,
	id int,
	limit core_domain.PatchLimit,
) (core_domain.Limit, error) {
	if err := limit.Validate(); err != nil {
		return core_domain.Limit{}, fmt.Errorf("limit validation failed: %w", err)
	}

	limitDomain, err := s.GetLimit(ctx, id)
	if err != nil {
		return core_domain.Limit{}, fmt.Errorf("limit get failed: %w", err)
	}

	if err := limitDomain.ApplyPatch(limit); err != nil {
		return core_domain.Limit{}, fmt.Errorf("limit apply patch error: %w", err)
	}
	updateLimit, err := s.limitRepo.PatchLimit(ctx, limitDomain)
	if err != nil {
		return core_domain.Limit{}, fmt.Errorf("limit update failed: %w", err)
	}

	return updateLimit, nil
}
