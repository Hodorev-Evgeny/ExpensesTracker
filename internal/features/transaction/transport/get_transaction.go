package feature_transactio_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// GetTransaction		godoc
// @Summary 			Get transaction
// @Description 		Get transaction by id
// @Tags 				transactions
// @Accept 				json
// @Produce 			json
// @Param				id				path int true "ID transaction"
// @Success				200	{object}	TransactionResponse "Get transaction successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure 			404	{object}	response.ErrorResponse "Not Found"
// @Router 				/transactions/{id}	[get]
func (h *TransactionHTTPHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	transactionID, err := core_http_query_parm.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "failed getting transaction id")
		return
	}

	transaction, err := h.transactionService.GetTransaction(ctx, transactionID)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "failed getting transaction")
		return
	}

	TransactionResp := ToTransactionResponse(transaction)

	ResponseHandler.JSONResponseHandler(http.StatusOK, TransactionResp)
}
