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
	Title core_http_types.Nullable[string] `json:"title" validate:"required"`
}

func (u *CategoryUpdateRequest) Validate() error {
	if u.Title.Set && u.Title.Value != nil {
		if len(*u.Title.Value) < 3 || len(*u.Title.Value) > 30 {
			return fmt.Errorf("title must be between 3 and 30 characters long")
		}
	}

	return nil
}

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

	ResponseHandler.JSONResponseHandler(http.StatusOK, categoryDomain)
}
