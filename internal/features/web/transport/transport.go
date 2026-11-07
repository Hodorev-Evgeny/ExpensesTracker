package webtransport

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	httpserver "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
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

func (h *WebTransport) Router() []httpserver.Route {
	return []httpserver.Route{
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
