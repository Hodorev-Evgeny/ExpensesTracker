//go:build integration

package feature_repository_transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_testutil "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/testutil"
)

func TestIntegrationTransactionRepository(t *testing.T) {
	ctx := context.Background()

	if err := TestPool.Ping(ctx); err != nil {
		t.Error(err)
	}

	testTransactionRepository := NewTransactionRepository(TestPool)

	now := time.Now().UTC()
	userA, err := core_testutil.CreateTestUser(ctx, TestPool, core_domain.User{
		Full_name: "John Smith",
		Email:     "testA@gmail.com",
		Password:  "12345789",
		Time_add:  now,
	})
	if err != nil {
		t.Fatalf("failed to create user A: %v", err)
	}

	userB, err := core_testutil.CreateTestUser(ctx, TestPool, core_domain.User{
		Full_name: "Patrick Jane",
		Email:     "testB@gmail.com",
		Password:  "12345789",
		Time_add:  now,
	})
	if err != nil {
		t.Fatalf("failed to create user B: %v", err)
	}

	categoryA, err := core_testutil.CreateTestCategory(ctx, TestPool, core_domain.Category{
		Name:    "Food",
		User_ID: userA.ID,
	})
	if err != nil {
		t.Fatalf("failed to create category A: %v", err)
	}

	transactionA, err := core_testutil.CreateTestTransaction(ctx, TestPool, core_domain.Transaction{
		Sum:         1500,
		Type:        "Expenditure",
		Date:        now,
		CategoryID:  categoryA.ID,
		UserID:      userA.ID,
		Comments:    "Lunch at restaurant",
		TimeCreated: now,
	})
	if err != nil {
		t.Fatalf("failed to create transaction A: %v", err)
	}

	t.Run("B cannot get transaction A", func(t *testing.T) {

		_, err := testTransactionRepository.GetTransaction(ctx, transactionA.ID, userB.ID)

		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Fatalf(
				"expected ErrorNotFoud, got: %v",
				err,
			)
		}
	})

	t.Run("B cannot patch transaction A", func(t *testing.T) {
		_, err := testTransactionRepository.PatchTransaction(ctx, transactionA, userB.ID)

		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Fatalf(
				"expected ErrorNotFoud, got: %v",
				err,
			)
		}
	})

	t.Run("B cannot delete transaction A", func(t *testing.T) {
		err := testTransactionRepository.DeleteTransaction(ctx, transactionA.ID, userB.ID)

		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Fatalf(
				"expected ErrorNotFoud, got: %v",
				err,
			)
		}
	})
}
