package feature_transactio_transport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type TransactionService interface {
	CreateTransaction(
		ctx context.Context,
		transaction core_domain.Transaction,
	) (core_domain.Transaction, error)
}

type TransactionHTTPHandler struct {
	transactionService TransactionService
}

func NewTransactionHTTPHandler(
	transactionService TransactionService,
) *TransactionHTTPHandler {
	return &TransactionHTTPHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHTTPHandler) Router() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/transaction",
			Handler: h.CreateTransaction,
		},
	}
}
