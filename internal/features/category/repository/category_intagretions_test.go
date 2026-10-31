//go:build integration

package feature_repositor_category

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_testutil "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/testutil"
)

func TestIntegrationCategory(t *testing.T) {
	ctx := context.Background()

	if err := TestPool.Ping(ctx); err != nil {
		t.Error(err)
	}

	testCategoryRepository := NewCategoryRepository(TestPool)

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

	t.Run("UserB cant get category userA", func(t *testing.T) {
		_, err := testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userB.ID)
		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrorNotFoud,
				err,
			)
		}
	})

	t.Run("UserB cant patch category userA", func(t *testing.T) {
		_, err := testCategoryRepository.RenameCategory(ctx, categoryA.ID, userB.ID, categoryA)

		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrorNotFoud,
				err,
			)
		}
	})

	t.Run("UserB cant delete category userA", func(t *testing.T) {
		err := testCategoryRepository.DeleteCategory(ctx, categoryA.ID, userB.ID)

		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrorNotFoud,
				err,
			)
		}
	})
}
