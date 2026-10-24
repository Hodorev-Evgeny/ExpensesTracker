package feature_repository_transaction

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	transaction_repository_mock "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/repository/mock"
)

type TestData struct {
	name          string
	pool          core_repository_pool.Pool
	function      func(ctx context.Context, repo *TransactionRepository) (any, error)
	expectedError error
}

func TestExtraditionUser(t *testing.T) {
	testPool := transaction_repository_mock.CreateTestPool()
	repo := NewTransactionRepository(testPool)
	testData := []TestData{
		{
			name: "Test get user",
			pool: testPool,
			function: func(ctx context.Context, repo *TransactionRepository) (any, error) {
				filters := core_domain.FiltersTransaction{
					Limit:  new(5),
					Offset: new(5),
				}
				return repo.GetTransactions(
					ctx,
					filters,
				)
			},
			expectedError: core_errors.ErrorMock,
		},
	}

	for _, test := range testData {
		_, err := test.function(context.Background(), repo)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Test %s expected error %v, but got %v", test.name, test.expectedError, err)
		}
	}
}
