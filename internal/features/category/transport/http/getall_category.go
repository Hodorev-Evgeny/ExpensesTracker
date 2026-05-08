package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

type CategoryListResponse []CategoryResponse

// GetAllCategorys		godoc
// @Summary 			Get all category
// @Description 		Get all category without query parm
// @Tags 				category
// @Accept 				json
// @Produce 			json
// @Success				200	{array}		CategoryResponse "Get all category successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Router 				/category		[get]
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
