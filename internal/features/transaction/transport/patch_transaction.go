package feature_transactio_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_types "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/types"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type TransactionPatchRequest struct {
	Sum             core_http_types.Nullable[int]    `json:"sum" swaggertype:"integer" example:"1023"`
	CategoryID      core_http_types.Nullable[int]    `json:"categoryId" swaggertype:"integer" example:"1"`
	TypeTransaction core_http_types.Nullable[string] `json:"typeTransaction" swaggertype:"string" example:"Expenditure/Income"`
	Comments        core_http_types.Nullable[string] `json:"comments" swaggertype:"string" example:"comment"`
}

// PatchTransaction		godoc
// @Summary 			Patch transaction
// @Description 		Patch transaction you can get all parm of nothing
// @Tags 				transactions
// @Accept 				json
// @Produce 			json
// @Param				id				path int true "ID transaction"
// @Param				request body	TransactionPatchRequest true "transaction filters for get"
// @Success				200	{object}	TransactionResponse "Patch transaction successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure 			404	{object}	response.ErrorResponse "Not Found"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router 				/transactions/{id}	[patch]
func (h *TransactionHTTPHandler) PatchTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	transactionId, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "failed to fetch transaction id")
		return
	}

	var request TransactionPatchRequest
	if err := core_http_utils.DecodeJSON(&request, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding request body")
		return
	}

	patchDomain := ToDomainPatch(request)
	updateTransaction, err := h.transactionService.PatchTransaction(ctx, transactionId, patchDomain)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error patching transaction")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusOK, updateTransaction)
}
