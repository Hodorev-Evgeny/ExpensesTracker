package statstransport

import (
	"context"
	"net/http"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	httpserver "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
)

type StaticService interface {
	GetStatic(
		ctx context.Context,
		filters core_domain.FiltersStatic,
	) (core_domain.Static, error)

	GetBudgetStatus(
		ctx context.Context,
		user_id int64,
		category_id core_domain.Nullable[int64],
		at core_domain.Nullable[time.Time],
	) (*pb.ListBudgetStatusesResponse, error)
}

type StaticHTTPHandler struct {
	StaticService StaticService
}

func NewStaticHTTPHandler(
	s StaticService,
) *StaticHTTPHandler {
	return &StaticHTTPHandler{
		StaticService: s,
	}
}

func (h *StaticHTTPHandler) Router() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodGet,
			Path:    "/static",
			Handler: h.GetStatic,
		},
	}
}
