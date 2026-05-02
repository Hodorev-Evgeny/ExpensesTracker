package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *CategoryHTTPHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	categoryID, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting category id")
		return
	}

	categoryDomain, err := h.CategoryService.GetCategoryByID(ctx, categoryID)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting category id")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusOK, categoryDomain)
}
