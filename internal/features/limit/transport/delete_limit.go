package feature_transport_limit

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// DeleteLimit			godoc
// @Summary 			Delete Limit
// @Description 		Delete limit in database
// @Tags 				limit
// @Success				204	"Delete limit successfully"
// @Param				id 				path int true "ID limit"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure 			404	{object}	response.ErrorResponse "Not found"
// @Router 				/limit/{id}			[delete]
func (h *LimitHTTPHandler) DeleteLimit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	limitID, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting limit id")
		return
	}

	if err := h.limitService.DeleteLimit(ctx, limitID); err != nil {
		ResponseHandler.ErrorResponse(err, "error deleting limit")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusNoContent, "")
}
