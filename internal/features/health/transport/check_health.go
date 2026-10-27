package feature_transport_health

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

type DataHealthResponse struct {
	Status string `json:"status"`
}

func (h *HealthTransport) CheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	errorHealth := h.repo.CheckHealth(ctx)
	if errorHealth != nil {
		ResponseHandler.ErrorResponse(errorHealth, "health check failed")
		return
	}

	rsp := DataHealthResponse{
		Status: "health check passed",
	}
	ResponseHandler.JSONResponseHandler(http.StatusOK, rsp)
}
