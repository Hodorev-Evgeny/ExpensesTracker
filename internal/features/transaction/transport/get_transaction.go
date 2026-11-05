package feature_transactio_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
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
// @Failure      		500 {object} 	response.ErrorResponse "Internal server error"
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

	userID, ok := core_domain.UserIDFromContext(ctx)
	if !ok {
		ResponseHandler.ErrorResponse(
			fmt.Errorf("user not found in context: %w", core_errors.ErrorUnauthorized),
			"user is not registered")
		return
	}

	transaction, err := h.transactionService.GetTransaction(ctx, transactionID, userID)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "failed getting transaction")
		return
	}

	TransactionResp := ToTransactionResponse(transaction)

	ResponseHandler.JSONResponseHandler(http.StatusOK, TransactionResp)
}
