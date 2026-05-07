package feature_transport_limit

import (
	"net/http"
	"time"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_types "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/types"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type PatchLimit struct {
	Duraction   core_http_types.Nullable[time.Time] `json:"duration"`
	AmountLimit core_http_types.Nullable[int]       `json:"amount_limit"`
}

func (h *LimitHTTPHandler) PatchLimit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	var patchLimit PatchLimit
	if err := core_http_utils.DecodeJSON(&patchLimit, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding patch limit")
		return
	}

	id, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding id patch limit")
		return
	}

	limitDomain := LimitResponseToPatch(patchLimit)

	updateLimit, err := h.limitService.PatchLimit(ctx, id, limitDomain)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error patch limit")
		return
	}

	resp := LimitDomainToResponse(updateLimit)
	ResponseHandler.JSONResponseHandler(http.StatusOK, resp)
}
