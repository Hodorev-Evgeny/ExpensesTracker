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
		id int,
		limit core_domain.PatchLimit,
	) (core_domain.Limit, error)

	DeleteLimit(
		ctx context.Context,
		limitID int,
	) error
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
		{
			Method:  http.MethodGet,
			Path:    "/limit/{id}",
			Handler: h.GetLimit,
		},
		{
			Method:  http.MethodGet,
			Path:    "/limits",
			Handler: h.GetLimits,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/limit/{id}",
			Handler: h.PatchLimit,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/limit/{id}",
			Handler: h.DeleteLimit,
		},
	}
}
