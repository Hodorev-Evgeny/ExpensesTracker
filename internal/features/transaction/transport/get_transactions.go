package feature_transactio_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// GetsTransaction		godoc
// @Summary 			Get all transaction
// @Description 		Get all transaction with filters you can get all parm of nothing
// @Tags 				transactions
// @Accept 				json
// @Produce 			json
// @Param				user_id 		query int false "Filter by user id"
// @Param				category_id 	query int false "Filter by user category id"
// @Param				sum 			query int false "Filter by more sum transaction"
// @Param				to 				query int false "Filter by time to"
// @Param				from 			query int false "Filter by time from"
// @Param				limit 			query int false "Filter for pagination"
// @Param				offset 			query int false "Filter for pagination"
// @Success				200	{object}	TransactionResponse "Get transaction successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router 				/transactions	[get]
func (h *TransactionHTTPHandler) GetsTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	filters, err := GetTransactionFilters(r)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting transaction filters")
		return
	}

	list, err := h.transactionService.GetTransactions(ctx, filters)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting transactions")
		return
	}

	listResponse := ListDomainFromResponse(list)
	ResponseHandler.JSONResponseHandler(http.StatusOK, listResponse)
}

func ListDomainFromResponse(
	transaction []core_domain.Transaction,
) []TransactionResponse {
	tmp := make([]TransactionResponse, len(transaction))
	for i, tx := range transaction {
		tmp[i] = ToTransactionResponse(tx)
	}
	return tmp
}

func GetTransactionFilters(r *http.Request) (core_domain.FiltersTransaction, error) {
	var filters core_domain.FiltersTransaction

	intFields := map[string]**int{
		"user_id":     &filters.UserId,
		"category_id": &filters.CategoryId,
		"sum":         &filters.Sum,
		"limit":       &filters.Limit,
		"offset":      &filters.Offset,
	}

	for key, field := range intFields {
		val, err := core_http_utils.GetIntQueryParm(r, key)
		if err != nil {
			return core_domain.FiltersTransaction{}, fmt.Errorf("error getting query param %q: %w", key, err)
		}

		if val != nil {
			*field = val
		}
	}

	to, err := core_http_utils.GetDateQueryParm(r, "to")
	if err != nil {
		return core_domain.FiltersTransaction{}, fmt.Errorf("error getting query param: %w", err)
	}
	from, err := core_http_utils.GetDateQueryParm(r, "from")
	if err != nil {
		return core_domain.FiltersTransaction{}, fmt.Errorf("error getting query param: %w", err)
	}

	filters.From = from
	filters.To = to

	return filters, nil
}
