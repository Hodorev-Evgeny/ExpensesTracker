package feature_category_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_types "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/types"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type CategoryUpdateRequest struct {
	Title core_http_types.Nullable[string] `json:"title" validate:"required" swaggertype:"string"`
}

func (u *CategoryUpdateRequest) Validate() error {
	if u.Title.Set && u.Title.Value != nil {
		if len(*u.Title.Value) < 3 || len(*u.Title.Value) > 30 {
			return fmt.Errorf("title must be between 3 and 30 characters long")
		}
	}

	return nil
}

// RenameCategory		godoc
// @Summary 			Patch category
// @Description 		Patch category in database by id you can give all param or nothing
// @Tags 				category
// @Accept 				json
// @Produce 			json
// @Param 				id	path int	true	"Id category"
// @Param				request body 			CategoryUpdateRequest true  "Parm from path category"
// @Success				200	{object}			CategoryResponse "Patch category successfully"
// @Failure 			400	{object}			response.ErrorResponse "Bad request"
// @Failure 			404	{object}			response.ErrorResponse "Category not found"
// @Failure      		500 {object} 			response.ErrorResponse "Internal server error"
// @Router 				/category/{id}			[patch]
func (h *CategoryHTTPHandler) RenameCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	categoryID, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "id is required")
		return
	}

	var request CategoryUpdateRequest
	if err := core_http_utils.DecodeJSON(&request, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding request")
		return
	}

	categoryUpdate := core_domain.RequestUpdateFromDomain(request.Title.ToDomain())

	categoryDomain, err := h.CategoryService.RenameCategory(ctx, categoryID, categoryUpdate)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error renaming category")
		return
	}

	resp := DomainFromResponse(categoryDomain)

	ResponseHandler.JSONResponseHandler(http.StatusOK, resp)
}
