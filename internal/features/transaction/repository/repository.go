package feature_repository_transaction

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type TransactionRepository struct {
	pool core_repository_pool.Pool
}

func NewTransactionRepository(
	pool core_repository_pool.Pool,
) *TransactionRepository {
	return &TransactionRepository{
		pool: pool,
	}
}
