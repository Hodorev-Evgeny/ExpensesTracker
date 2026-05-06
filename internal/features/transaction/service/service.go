package feature_transaction_service

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type TransactionRepository interface {
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
		transaction core_domain.Transaction,
	) (core_domain.Transaction, error)

	DeleteTransaction(
		ctx context.Context,
		id int,
	) error
}

type TransactionService struct {
	transactionRepository TransactionRepository
}

func NewTransactionService(
	transactionRepository TransactionRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}
