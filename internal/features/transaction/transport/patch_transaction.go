package feature_transactio_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_types "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/types"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type TransactionPatchRequest struct {
	Sum             core_http_types.Nullable[int]    `json:"sum"`
	CategoryID      core_http_types.Nullable[int]    `json:"categoryId"`
	TypeTransaction core_http_types.Nullable[string] `json:"typeTransaction"`
	Comments        core_http_types.Nullable[string] `json:"comments"`
}

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
