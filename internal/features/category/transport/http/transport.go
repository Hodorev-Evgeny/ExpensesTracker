package feature_category_transport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type CategoryService interface {
	CreateCategory(
		ctx context.Context,
		category core_domain.Category,
	) (core_domain.Category, error)
}

type CategoryHTTPHandler struct {
	CategoryService CategoryService
}

func NewCategoryHTTPHandler(
	categoryService CategoryService,
) *CategoryHTTPHandler {
	return &CategoryHTTPHandler{
		CategoryService: categoryService,
	}
}

func (h *CategoryHTTPHandler) Routes() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/category",
			Handler: h.CreateNewCategory,
		},
	}
}
