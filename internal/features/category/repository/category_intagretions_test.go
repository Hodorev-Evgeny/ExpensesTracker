//go:build integration

package categoryrepo

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

	if _, err = core_testutil.CreateTestCategory(ctx, TestPool, core_domain.Category{
		Name:    "Bar",
		User_ID: userB.ID,
	}); err != nil {
		t.Fatalf("failed to create category B: %v", err)
	}

	t.Run("UserB cant get category userA", func(t *testing.T) {
		_, err := testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userB.ID)
		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrNotFound,
				err,
			)
		}
	})

	t.Run("UserB cant patch category userA", func(t *testing.T) {
		patch := categoryA
		patch.Name = "New category"

		_, err := testCategoryRepository.RenameCategory(ctx, categoryA.ID, userB.ID, patch)

		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrNotFound,
				err,
			)
		}

		ans, err := testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userA.ID)
		if err != nil {
			t.Errorf("Wanted no error, got: %v",
				err,
			)
		}

		assert.Equal(t, ans, categoryA)
	})

	t.Run("UserB cant delete category userA", func(t *testing.T) {
		err := testCategoryRepository.DeleteCategory(ctx, categoryA.ID, userB.ID)

		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrNotFound,
				err,
			)
		}

		ans, err := testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userA.ID)
		if err != nil {
			t.Errorf("Wanted no error, got: %v",
				err,
			)
		}

		assert.Equal(t, ans, categoryA)
	})

	t.Run("UserA can see only their category", func(t *testing.T) {
		list, err := testCategoryRepository.GetAllCategories(ctx, userA.ID)
		if err != nil {
			t.Errorf("Wanted no error, got: %v",
				err,
			)
		}

		flag := false
		for _, category := range list {
			if category.User_ID != userA.ID {
				t.Errorf("got user ID %v, wanted user ID %v", category.User_ID, userA.ID)
			}
			if category.ID == categoryA.ID {
				flag = true
			}
		}

		if !flag {
			t.Error("In list not found category userA")
		}
	})

	t.Run("UserA can get their transaction", func(t *testing.T) {
		ans, err := testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userA.ID)
		if err != nil {
			t.Errorf("Wanted no error, got: %v",
				err,
			)
		}

		assert.Equal(t, ans, categoryA)
	})

	t.Run("UserA ca patch their transaction", func(t *testing.T) {
		patch := categoryA
		patch.Name = "New category"

		_, err := testCategoryRepository.RenameCategory(ctx, categoryA.ID, userA.ID, patch)
		if err != nil {
			t.Errorf("Wanted no error, got: %v",
				err,
			)
		}

		ans, err := testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userA.ID)
		assert.Equal(t, ans, patch)
	})

	t.Run("UserA can delete their transaction", func(t *testing.T) {
		err := testCategoryRepository.DeleteCategory(ctx, categoryA.ID, userA.ID)
		if err != nil {
			t.Errorf("Wanted no error, got: %v",
				err,
			)
		}

		_, err = testCategoryRepository.GetCategoryByID(ctx, categoryA.ID, userA.ID)
		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Errorf("Wanted error: %v, got: %v",
				core_errors.ErrNotFound,
				err,
			)
		}
	})
}
