package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// GetCategory		godoc
// @Summary 		Get category
// @Description 	Get category in database by id
// @Tags 			category
// @Accept 			json
// @Produce 		json
// @Param 			id	path int		true	"Id category"
// @Success			200	{object}		CategoryResponse "Get category successfully"
// @Failure 		400	{object}		response.ErrorResponse "Bad request"
// @Failure 		404	{object}		response.ErrorResponse "Category not found"
// @Router 			/category/{id}		[get]
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
