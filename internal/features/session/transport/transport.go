package feature_transport_session

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type SessionService interface {
	GetSessionData(
		ctx context.Context,
		key string,
	) (core_domain.CookieData, error)
}

type SessionHandler struct {
	SessionService SessionService
}

func NewSessionHandler(
	s SessionService,
) *SessionHandler {
	return &SessionHandler{
		SessionService: s,
	}
}

func (h *SessionHandler) Route() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/session/{sessionID}",
			Handler: h.GetSessionData,
		},
	}
}
