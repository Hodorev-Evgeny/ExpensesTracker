package feature_category_transport

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type CategoryRequest struct {
	CategoryName string `json:"category_name" validate:"required,min=3,max=20"`
	User_id      int    `json:"user_id" validate:"required"`
}

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

func CategoryDTOFromDomain(category CategoryRequest) core_domain.Category {
	return core_domain.CreateUnincelizedCategory(
		category.CategoryName,
		category.User_id,
	)
}
