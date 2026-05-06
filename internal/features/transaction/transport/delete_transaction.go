package feature_transactio_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *TransactionHTTPHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	transactionID, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting transaction id")
		return
	}

	if err := h.transactionService.DeleteTransaction(ctx, transactionID); err != nil {
		ResponseHandler.ErrorResponse(err, "error deleting transaction")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusNoContent, "")
}
