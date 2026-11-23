package statsservice

import (
	"context"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
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

	GetBudgetStatus(
		ctx context.Context,
		user_id int64,
		category_id core_domain.Nullable[int64],
		at core_domain.Nullable[time.Time],
	) (*pb.ListBudgetStatusesResponse, error)
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
