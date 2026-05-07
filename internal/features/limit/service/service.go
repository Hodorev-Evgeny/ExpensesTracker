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

	GetLimit(
		ctx context.Context,
		id int,
	) (core_domain.Limit, error)

	GetLimits(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]core_domain.Limit, error)

	PatchLimit(
		ctx context.Context,
		limit core_domain.Limit,
	) (core_domain.Limit, error)

	DeleteLimit(
		ctx context.Context,
		limitID int,
	) error
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
