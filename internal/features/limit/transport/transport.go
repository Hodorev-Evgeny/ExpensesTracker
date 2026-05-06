package feature_transport_limit

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type LimitService interface {
	CreateLimit(
		ctx context.Context,
		limit core_domain.Limit,
	) (core_domain.Limit, error)
}

type LimitHTTPHandler struct {
	limitService LimitService
}

func NewLimitHTTPHandler(
	limitService LimitService,
) *LimitHTTPHandler {
	return &LimitHTTPHandler{
		limitService: limitService,
	}
}

func (h LimitHTTPHandler) Router() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/limit",
			Handler: h.CreateLimit,
		},
	}
}
