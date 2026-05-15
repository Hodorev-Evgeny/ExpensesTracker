package feature_transport_web

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type WebService interface {
	GetMainPage() (core_domain.File, error)
	GetRegister() (core_domain.File, error)
	GetLogin() (core_domain.File, error)
}

type WebTransport struct {
	service WebService
}

func NewWebTransport(
	service WebService,
) *WebTransport {
	return &WebTransport{
		service: service,
	}
}

func (h *WebTransport) Router() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{$}",
			Handler: h.GetMainPage,
		},
		{
			Method:  http.MethodGet,
			Path:    "/register",
			Handler: h.GetRegisterPage,
		},
		{
			Method:  http.MethodGet,
			Path:    "/login",
			Handler: h.GetLogin,
		},
	}
}
