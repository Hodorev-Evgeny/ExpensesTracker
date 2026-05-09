package feature_transactio_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// DeleteTransaction	godoc
// @Summary 			Delete transaction
// @Description 		Delete transaction in database
// @Tags 				transactions
// @Param				id				path int true "ID transaction"
// @Success				204				"Delete transaction successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure 			404	{object}	response.ErrorResponse "Not Found"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router 				/transactions/{id}	[delete]
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
