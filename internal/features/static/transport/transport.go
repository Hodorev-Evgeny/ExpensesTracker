package feature_transport_static

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type StaticService interface {
	GetStatic(
		ctx context.Context,
		filters core_domain.FiltersStatic,
	) (core_domain.Static, error)
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

func (h *StaticHTTPHandler) Router() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/static",
			Handler: h.GetStatic,
		},
	}
}
