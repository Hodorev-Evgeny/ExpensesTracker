package feature_category_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type CategoryRequest CategoryDTO

// CreateCategory 	godoc
// @Summary 		Create category
// @Description 	Create new category in database
// @Tags 			category
// @Accept 			json
// @Produce 		json
// @Param 			request body	CategoryRequest true "CreateCategory body"
// @Success			201	{object}	CategoryResponse "Create new category successfully"
// @Failure 		400	{object}	response.ErrorResponse "Bad request"
// @Router 			/category		[post]
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
