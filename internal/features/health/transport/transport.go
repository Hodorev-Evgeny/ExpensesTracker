package feature_transport_health

import (
	"context"
	"net/http"

	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
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

func (h HealthTransport) Router() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: h.CheckHealth,
		},
	}
}
