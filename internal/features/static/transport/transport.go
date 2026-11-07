package statstransport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	httpserver "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
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

func (h *StaticHTTPHandler) Router() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodGet,
			Path:    "/static",
			Handler: h.GetStatic,
		},
	}
}
