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

	GetTransactions(
		ctx context.Context,
		filters core_domain.FiltersTransaction,
	) ([]core_domain.Transaction, error)

	GetTransaction(
		ctx context.Context,
		id int,
	) (core_domain.Transaction, error)

	PatchTransaction(
		ctx context.Context,
		id int,
		patch core_domain.PatchTransaction,
	) (core_domain.Transaction, error)

	DeleteTransaction(
		ctx context.Context,
		id int,
	) error
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
		{
			Method:  http.MethodGet,
			Path:    "/transactions",
			Handler: h.GetsTransaction,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions/{id}",
			Handler: h.GetTransaction,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/transactions/{id}",
			Handler: h.PatchTransaction,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/transactions/{id}",
			Handler: h.DeleteTransaction,
		},
	}
}
