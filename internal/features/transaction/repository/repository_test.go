package feature_repository_transaction

import (
	"context"
	"regexp"
	"strconv"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	transaction_repository_mock "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/repository/mock"
)

type TestData struct {
	name           string
	pool           *transaction_repository_mock.TestPool
	function       func(ctx context.Context, repo *TransactionRepository) (any, error)
	expectedResult [2]any // first is limit second is offset
}

func TestGetTransactionsLimitOffset(t *testing.T) {
	testPool := transaction_repository_mock.CreateTestPool()
	repo := NewTransactionRepository(testPool)
	testData := []TestData{
		{
			name: "Test get correct limit and offset",
			pool: testPool,
			function: func(ctx context.Context, repo *TransactionRepository) (any, error) {
				filters := core_domain.FiltersTransaction{
					Limit:  new(5),
					Offset: new(7),
				}
				return repo.GetTransactions(
					ctx,
					filters,
				)
			},
			expectedResult: [2]any{5, 7},
		},
		{
			name: "Test correct read args",
			pool: testPool,
			function: func(ctx context.Context, repo *TransactionRepository) (any, error) {
				filters := core_domain.FiltersTransaction{
					Limit:      new(2),
					Offset:     new(3),
					CategoryId: new(1),
				}
				return repo.GetTransactions(
					ctx,
					filters,
				)
			},
			expectedResult: [2]any{2, 3},
		},
	}

	for _, test := range testData {
		_, err := test.function(context.Background(), repo)
		if err != nil {
			t.Error("expected error to be nil")
		}
		query, args := test.pool.GetQueryAndArgs()

		reLimit := regexp.MustCompile(`LIMIT\s+\$(\d+)`)
		matchLimit := reLimit.FindStringSubmatch(query)

		reOffset := regexp.MustCompile(`OFFSET\s+\$(\d+)`)
		matchOffset := reOffset.FindStringSubmatch(query)

		if len(matchLimit) > 1 {
			idx, err := strconv.Atoi(matchLimit[1])
			if err != nil {
				t.Error("Args idx for limit query is not a number")
				continue
			}

			if idx-1 < 0 || idx > len(args) {
				t.Fatalf("Args idx for limit query is out of range")
			}
			val := args[idx-1]
			if val != test.expectedResult[0] {
				t.Error("Args idx for limit query is not correct")
			}
		} else {
			t.Fatalf("Dont find limit query: %s", query)
		}

		if len(matchOffset) > 1 {
			idx, err := strconv.Atoi(matchOffset[1])
			if err != nil {
				t.Error("Args idx for offset query is not a number")
				continue
			}

			if idx-1 < 0 || idx > len(args) {
				t.Fatalf("Args idx for offset query is out of range")
			}
			val := args[idx-1]
			if val != test.expectedResult[1] {
				t.Error("Args idx for offset query is not correct")
			}
		} else {
			t.Fatalf("Dont find offset query: %s", query)
		}

	}
}
