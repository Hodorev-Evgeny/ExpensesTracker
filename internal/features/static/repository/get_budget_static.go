package statsrepo

import (
	"context"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
)

func (r *StaticRepository) GetBudgetStatus(
	ctx context.Context,
	user_id int64,
	category_id core_domain.Nullable[int64],
	at core_domain.Nullable[time.Time],
) (*pb.ListBudgetStatusesResponse, error) {
	args := grpcclient.OptionalArgsStatus{
		CategoryID: category_id,
		At:         at,
	}

	ans, err := r.client.ListBudgetStatuses(ctx, user_id, args)
	if err != nil {
		return nil, err
	}
	return ans, nil
}
