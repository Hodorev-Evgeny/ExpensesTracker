package feature_service_limit

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type LimitRepository interface {
	CreateLimit(
		ctx context.Context,
		limit core_domain.Limit,
	) (core_domain.Limit, error)
}

type LimitService struct {
	limitRepo LimitRepository
}

func NewLimitService(
	limitRepo LimitRepository,
) *LimitService {
	return &LimitService{
		limitRepo: limitRepo,
	}
}
