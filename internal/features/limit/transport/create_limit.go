package feature_transport_limit

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// CreateLimit			godoc
// @Summary 			Create new Limit
// @Description 		This limit dont set on category, you need patch category
// @Tags 				limit
// @Accept 				json
// @Produce 			json
// @Param				request body	LimitDTO true "Body for create limit"
// @Success				201	{object}	LimitResponse "Create new limit successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Router 				/limit			[post]
func (h *LimitHTTPHandler) CreateLimit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	var req LimitDTO
	if err := core_http_utils.DecodeJSON(&req, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding request body")
		return
	}

	limitUnincelizd := LimitResponseToDomain(req)

	limitDomain, err := h.limitService.CreateLimit(ctx, limitUnincelizd)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error creating limit")
		return
	}

	limitResp := LimitDomainToResponse(limitDomain)

	ResponseHandler.JSONResponseHandler(http.StatusCreated, limitResp)
}
