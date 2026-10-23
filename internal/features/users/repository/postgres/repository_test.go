package features_users_repository

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	mock_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/repository/postgres/mock"
)

type TestData struct {
	name           string
	pool           core_repository_pool.Pool
	function       func(ctx context.Context, repo *UserRepository) (any, error)
	expectedError  error
	expectedResult core_domain.User
}

func ExtraditionUserTest(t *testing.T) {
	testPool := mock_repository.CreateTestPool()
	repo := NewUserRepository(testPool)
	testData := []TestData{
		{
			name: "Test get user",
			pool: testPool,
			function: func(ctx context.Context, repo *UserRepository) (any, error) {
				return repo.ExtraditionUser(ctx, -1)
			},
			expectedError:  core_repository_pool.ErrNoRows,
			expectedResult: core_domain.User{},
		},
	}

	for _, test := range testData {
		ans, err := test.function(context.Background(), repo)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Test %s expected result %v, but got %v", test.name, test.expectedResult, err)
		}
		if ans != test.expectedResult {
			t.Errorf("Test %s expected result %v, but got %v", test.name, test.expectedResult, ans)
		}
	}
}
