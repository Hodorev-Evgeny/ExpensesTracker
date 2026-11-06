package healthtransport

import (
	"context"
	"net/http"

	httpserver "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type HealthRepository interface {
	CheckHealth(ctx context.Context) error
}

type HealthTransport struct {
	repo HealthRepository
}

func NewHealthTransport(
	repo HealthRepository,
) *HealthTransport {
	return &HealthTransport{repo: repo}
}

func (h HealthTransport) Router() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: h.CheckHealth,
		},
	}
}
