package limittransport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	httpserver "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type LimitService interface {
	CreateLimit(
		ctx context.Context,
		limit core_domain.Limit,
	) (core_domain.Limit, error)

	GetLimit(
		ctx context.Context,
		id int,
		userID int,
	) (core_domain.Limit, error)

	GetLimits(
		ctx context.Context,
		limit *int,
		offset *int,
		userID int,
	) ([]core_domain.Limit, error)

	PatchLimit(
		ctx context.Context,
		id int,
		limit core_domain.PatchLimit,
		userID int,
	) (core_domain.Limit, error)

	DeleteLimit(
		ctx context.Context,
		limitID int,
		userID int,
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

func (h LimitHTTPHandler) Router() []httpserver.Route {
	return []httpserver.Route{
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
			Path:    "/limit",
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
