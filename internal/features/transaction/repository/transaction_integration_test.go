//go:build integration

package transactionrepo

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_testutil "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/testutil"
	"github.com/magiconair/properties/assert"
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

	categoryB, err := core_testutil.CreateTestCategory(ctx, TestPool, core_domain.Category{
		Name:    "pizza",
		User_ID: userB.ID,
	})
	if err != nil {
		t.Fatalf("failed to create category A: %v", err)
	}

	if _, err = core_testutil.CreateTestTransaction(ctx, TestPool, core_domain.Transaction{
		Sum:         15,
		Type:        "Expenditure",
		Date:        now,
		CategoryID:  categoryB.ID,
		UserID:      userB.ID,
		Comments:    "Lunch at restaurant",
		TimeCreated: now,
	}); err != nil {
		t.Fatalf("failed to create transaction B: %v", err)
	}

	t.Run("B cannot get transaction A", func(t *testing.T) {
		_, err := testTransactionRepository.GetTransaction(ctx, transactionA.ID, userB.ID)

		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Fatalf(
				"expected ErrNotFound, got: %v",
				err,
			)
		}
	})

	t.Run("B cannot patch transaction A", func(t *testing.T) {
		newTransaction := transactionA
		newTransaction.Sum = 12345

		_, err := testTransactionRepository.PatchTransaction(ctx, newTransaction, userB.ID)

		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Fatalf(
				"expected ErrNotFound, got: %v",
				err,
			)
		}

		ans, err := testTransactionRepository.GetTransaction(ctx, transactionA.ID, userA.ID)
		if err != nil {
			t.Fatalf("failed to get transaction A: %v", err)
		}
		assert.Equal(t, ans, transactionA)
	})

	t.Run("B cannot delete transaction A", func(t *testing.T) {
		err := testTransactionRepository.DeleteTransaction(ctx, transactionA.ID, userB.ID)

		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Fatalf(
				"expected ErrNotFound, got: %v",
				err,
			)
		}

		_, err = testTransactionRepository.GetTransaction(ctx, transactionA.ID, userA.ID)
		if err != nil {
			t.Fatalf("Wanted no error, got: %v",
				err,
			)
		}
	})

	t.Run("A see only their transaction", func(t *testing.T) {
		filter := core_domain.FiltersTransaction{
			UserId: &userA.ID,
		}

		list, err := testTransactionRepository.GetTransactions(ctx, filter)
		if err != nil {
			t.Fatalf("failed to get transactions: %v", err)
		}

		flag := false
		for _, transaction := range list {
			if transaction.UserID != userA.ID {
				t.Fatalf("got user ID %v, wanted user ID %v", transaction.UserID, userA.ID)
			}
			if transaction.ID == transactionA.ID {
				flag = true
			}
		}

		if !flag {
			t.Fatal("Transaction userA not found in list")
		}
	})

	t.Run("A correct get their transaction", func(t *testing.T) {
		ans, err := testTransactionRepository.GetTransaction(ctx, transactionA.ID, userA.ID)
		if err != nil {
			t.Fatalf(
				"expected error to be nil, got: %v",
				err,
			)
		}
		assert.Equal(t, ans, transactionA)
	})

	t.Run("A correct patch their transaction", func(t *testing.T) {
		newTransaction := transactionA
		newTransaction.Sum = 12345

		ans, err := testTransactionRepository.PatchTransaction(ctx, newTransaction, userA.ID)
		if err != nil {
			t.Fatalf(
				"expected error to be nil, got: %v",
				err,
			)
		}
		assert.Equal(t, ans, newTransaction)
	})

	t.Run("A correct delete their transaction", func(t *testing.T) {
		err := testTransactionRepository.DeleteTransaction(ctx, transactionA.ID, userA.ID)
		if err != nil {
			t.Fatalf("failed to delete transaction A: %v",
				err,
			)
		}

		_, err = testTransactionRepository.GetTransaction(ctx, transactionA.ID, userA.ID)

		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Fatalf(
				"expected ErrNotFound, got: %v",
				err,
			)
		}
	})
}
