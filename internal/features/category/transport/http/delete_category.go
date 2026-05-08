package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// DeleteCategory 	godoc
// @Summary 		Delete category
// @Description 	Delete category in database by id
// @Tags 			category
// @Param			id	path int true	"Id delete category"
// @Success			204					"Delete category by id successfully"
// @Failure 		400	{object}		response.ErrorResponse "Bad request"
// @Router 			/category/{id}		[delete]
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
