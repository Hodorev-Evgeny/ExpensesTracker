package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type CategoryRequest CategoryDTO

func (h *CategoryHTTPHandler) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	var req CategoryRequest
	if err := core_http_utils.DecodeJSON(&req, r); err != nil {
		ResponseHandler.ErrorResponse(err, "Error decoding request body")
		return
	}

	categoryDomain := CategoryDTOFromDomain(req)

	category, err := h.CategoryService.CreateCategory(ctx, categoryDomain)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "Error creating category")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusCreated, category)
}
