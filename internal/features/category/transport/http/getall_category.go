package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

type CategoryListResponse []CategoryResponse

func (h *CategoryHTTPHandler) GetAllCategorys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	listCategory, err := h.CategoryService.GetAllCategories(ctx)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting all categories")
	}

	ResponseHandler.JSONResponseHandler(http.StatusOK, listCategory)
}
