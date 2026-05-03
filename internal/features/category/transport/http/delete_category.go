package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *CategoryHTTPHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	id, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "Error getting ID")
		return
	}

	if err := h.CategoryService.DeleteCategory(ctx, id); err != nil {
		ResponseHandler.ErrorResponse(err, "Error deleting Category")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusNoContent, nil)
}
