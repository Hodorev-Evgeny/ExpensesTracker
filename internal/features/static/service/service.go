package feature_service_static

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type StaticeRepository interface {
	GetStaticCategories(
		ctx context.Context,
		filters core_domain.FiltersStatic,
	) ([]core_domain.Category, error)

	GetStaticTransactions(
		ctx context.Context,
		filters core_domain.FiltersStatic,
	) ([]core_domain.Transaction, error)
}

type StaticService struct {
	staticRepository StaticeRepository
}

func NewStaticService(
	repository StaticeRepository,
) *StaticService {
	return &StaticService{
		staticRepository: repository,
	}
}
