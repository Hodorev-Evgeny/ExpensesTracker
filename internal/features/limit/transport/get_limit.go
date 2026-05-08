package feature_transport_limit

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// GetLimit				godoc
// @Summary 			Get Limit by id
// @Description 		Get limit in database ny id
// @Tags 				limit
// @Success				200				"Get limit successfully"
// @Param				id 				path int true "ID limit"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure 			404	{object}	response.ErrorResponse "Not found"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router 				/limit/{id}		[get]
func (h *LimitHTTPHandler) GetLimit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	id, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "decoding limit id")
		return
	}

	limitDomain, err := h.limitService.GetLimit(ctx, id)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "getting limit domain")
		return
	}

	req := LimitDomainToResponse(limitDomain)
	ResponseHandler.JSONResponseHandler(http.StatusOK, req)
}
