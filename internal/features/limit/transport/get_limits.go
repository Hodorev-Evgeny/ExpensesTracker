package feature_transport_limit

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *LimitHTTPHandler) GetLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	limit, err := core_http_utils.GetIntQueryParm(r, "limit")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting limit")
		return
	}
	offset, err := core_http_utils.GetIntQueryParm(r, "offset")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting offset")
		return
	}

	limitDomain, err := h.limitService.GetLimits(ctx, limit, offset)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "getting limit domain")
		return
	}

	req := ListToResp(limitDomain)
	ResponseHandler.JSONResponseHandler(http.StatusOK, req)
}

func ListToResp(l []core_domain.Limit) []LimitResponse {
	resp := make([]LimitResponse, len(l))
	for i := range l {
		resp[i] = LimitDomainToResponse(l[i])
	}
	return resp
}
