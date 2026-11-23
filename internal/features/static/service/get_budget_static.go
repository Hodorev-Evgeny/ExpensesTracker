package statsservice

import (
	"context"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
)

func (s *StaticService) GetBudgetStatus(
	ctx context.Context,
	user_id int64,
	category_id core_domain.Nullable[int64],
	at core_domain.Nullable[time.Time],
) (*pb.ListBudgetStatusesResponse, error) {
	return s.staticRepository.GetBudgetStatus(ctx, user_id, category_id, at)
}
