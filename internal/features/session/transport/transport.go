package feature_transport_session

import (
	"net/http"

	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type SessionHandler struct {
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}

func (h *SessionHandler) Route() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users/me",
			Handler: h.GetSessionData,
		},
	}
}
