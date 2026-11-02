package feature_service_categor

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestServiceCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("Create correct category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		category := core_domain.Category{
			ID:       1,
			Name:     "test",
			User_ID:  12,
			Limit_id: nil,
		}

		repo.EXPECT().
			CreateCategory(ctx, category).
			Return(category, nil).
			Times(1)

		result, err := service.CreateCategory(ctx, category)

		assert.NoError(t, err)
		assert.Equal(t, category, result)
	})

	t.Run("Create invalid category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		category := core_domain.Category{
			ID:      1,
			Name:    "",
			User_ID: 12,
		}

		result, err := service.CreateCategory(ctx, category)

		assert.Error(t, err)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Create category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)
		repoErr := errors.New("repository create failed")

		category := core_domain.Category{
			ID:      1,
			Name:    "test",
			User_ID: 12,
		}

		repo.EXPECT().
			CreateCategory(ctx, category).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.CreateCategory(ctx, category)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Get existing category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		category := core_domain.Category{
			ID:      1,
			Name:    "test",
			User_ID: 12,
		}

		repo.EXPECT().
			GetCategoryByID(ctx, category.ID).
			Return(category, nil).
			Times(1)

		result, err := service.GetCategoryByID(ctx, category.ID)

		assert.NoError(t, err)
		assert.Equal(t, category, result)
	})

	t.Run("Get category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)
		repoErr := core_errors.ErrorNotFoud

		repo.EXPECT().
			GetCategoryByID(ctx, 1).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.GetCategoryByID(ctx, 1)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Get category invalid ID", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		result, err := service.GetCategoryByID(ctx, core_domain.UnincelizedID)

		assert.Error(t, err)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Successful rename category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		category := core_domain.Category{
			ID:      1,
			Name:    "test",
			User_ID: 12,
		}
		newName := "renamed"
		patch := core_domain.CategoryUpdate{
			Name: core_domain.Nullable[string]{
				Value: &newName,
				Set:   true,
			},
		}
		expected := category
		expected.Name = newName

		repo.EXPECT().
			GetCategoryByID(ctx, category.ID).
			Return(category, nil).
			Times(1)

		repo.EXPECT().
			RenameCategory(ctx, category.ID, expected).
			Return(expected, nil).
			Times(1)

		result, err := service.RenameCategory(ctx, category.ID, patch)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Rename category get error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)
		repoErr := core_errors.ErrorNotFoud

		newName := "renamed"
		patch := core_domain.CategoryUpdate{
			Name: core_domain.Nullable[string]{
				Value: &newName,
				Set:   true,
			},
		}

		repo.EXPECT().
			GetCategoryByID(ctx, 1).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.RenameCategory(ctx, 1, patch)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Rename category invalid patch", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		category := core_domain.Category{
			ID:      1,
			Name:    "test",
			User_ID: 12,
		}
		patch := core_domain.CategoryUpdate{
			Name: core_domain.Nullable[string]{
				Set: true,
			},
		}

		repo.EXPECT().
			GetCategoryByID(ctx, category.ID).
			Return(category, nil).
			Times(1)

		result, err := service.RenameCategory(ctx, category.ID, patch)

		assert.Error(t, err)
		assert.Equal(t, category, result)
	})

	t.Run("Rename category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)
		repoErr := errors.New("repository rename failed")

		category := core_domain.Category{
			ID:      1,
			Name:    "test",
			User_ID: 12,
		}
		newName := "renamed"
		patch := core_domain.CategoryUpdate{
			Name: core_domain.Nullable[string]{
				Value: &newName,
				Set:   true,
			},
		}
		updated := category
		updated.Name = newName

		repo.EXPECT().
			GetCategoryByID(ctx, category.ID).
			Return(category, nil).
			Times(1)

		repo.EXPECT().
			RenameCategory(ctx, category.ID, updated).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.RenameCategory(ctx, category.ID, patch)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Successful delete category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		repo.EXPECT().
			DeleteCategory(ctx, 1).
			Return(nil).
			Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("Delete category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)
		repoErr := errors.New("repository delete failed")

		repo.EXPECT().
			DeleteCategory(ctx, 1).
			Return(repoErr).
			Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("Delete category invalid ID", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		service := NewCategoryService(repo)

		err := service.DeleteCategory(ctx, core_domain.UnincelizedID)

		assert.Error(t, err)
	})
}
