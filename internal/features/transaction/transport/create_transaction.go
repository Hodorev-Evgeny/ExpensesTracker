package feature_transactio_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type TransactionRequest TransactionDTO

// CreateTransaction	godoc
// @Summary 			Create transaction
// @Description 		Create transaction in database
// @Tags 				transactions
// @Accept 				json
// @Produce 			json
// @Param				request body	TransactionRequest true "Example body request"
// @Success				201	{object}	TransactionResponse "Create transaction successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router 				/transactions	[post]
func (h *TransactionHTTPHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	var request TransactionRequest
	if err := core_http_utils.DecodeJSON(&request, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding request body")
		return
	}

	transactionDomain := ToDomainTransaction(request)

	transactionCreated, err := h.transactionService.CreateTransaction(ctx, transactionDomain)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error creating transaction")
		return
	}

	transactionResp := ToTransactionResponse(transactionCreated)

	ResponseHandler.JSONResponseHandler(http.StatusCreated, transactionResp)
}
