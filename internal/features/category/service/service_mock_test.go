package categoryservice

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	corelogger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type noopBudgetClient struct{}

func (noopBudgetClient) ApplySpending(
	ctx context.Context,
	transaction core_domain.Transaction,
) (*pb.ApplySpendingResponse, error) {
	return &pb.ApplySpendingResponse{}, nil
}

func (noopBudgetClient) RevokeSpending(
	ctx context.Context,
	transaction core_domain.DeletedTransaction,
) (*pb.RevokeSpendingResponse, error) {
	return &pb.RevokeSpendingResponse{}, nil
}

var _ grpcclient.BudgetService = noopBudgetClient{}

func TestServiceCategory(t *testing.T) {
	ctx := context.Background()
	ctx = corelogger.ToContext(ctx, &corelogger.Logger{
		Logger: zap.NewNop(),
	})

	t.Run("Create correct category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

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
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

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
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)
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
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

		category := core_domain.Category{
			ID:      1,
			Name:    "test",
			User_ID: 12,
		}

		repo.EXPECT().
			GetCategoryByID(ctx, category.ID, category.User_ID).
			Return(category, nil).
			Times(1)

		result, err := service.GetCategoryByID(ctx, category.ID, category.User_ID)

		assert.NoError(t, err)
		assert.Equal(t, category, result)
	})

	t.Run("Get category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)
		repoErr := core_errors.ErrNotFound

		repo.EXPECT().
			GetCategoryByID(ctx, 1, 12).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.GetCategoryByID(ctx, 1, 12)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Get category invalid ID", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

		result, err := service.GetCategoryByID(ctx, core_domain.UninitializedID, 12)

		assert.Error(t, err)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Successful rename category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

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
			GetCategoryByID(ctx, category.ID, category.User_ID).
			Return(category, nil).
			Times(1)

		repo.EXPECT().
			RenameCategory(ctx, category.ID, category.User_ID, expected).
			Return(expected, nil).
			Times(1)

		result, err := service.RenameCategory(ctx, category.ID, patch, category.User_ID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Rename category get error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)
		repoErr := core_errors.ErrNotFound

		newName := "renamed"
		patch := core_domain.CategoryUpdate{
			Name: core_domain.Nullable[string]{
				Value: &newName,
				Set:   true,
			},
		}

		repo.EXPECT().
			GetCategoryByID(ctx, 1, 12).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.RenameCategory(ctx, 1, patch, 12)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Rename category invalid patch", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

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
			GetCategoryByID(ctx, category.ID, category.User_ID).
			Return(category, nil).
			Times(1)

		result, err := service.RenameCategory(ctx, category.ID, patch, category.User_ID)

		assert.Error(t, err)
		assert.Equal(t, category, result)
	})

	t.Run("Rename category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)
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
			GetCategoryByID(ctx, category.ID, category.User_ID).
			Return(category, nil).
			Times(1)

		repo.EXPECT().
			RenameCategory(ctx, category.ID, category.User_ID, updated).
			Return(core_domain.Category{}, repoErr).
			Times(1)

		result, err := service.RenameCategory(ctx, category.ID, patch, category.User_ID)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Successful delete category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

		repo.EXPECT().
			DeleteCategory(ctx, 1, 12).
			Return([]core_domain.DeletedTransaction{}, nil).
			Times(1)

		err := service.DeleteCategory(ctx, 1, 12)

		assert.NoError(t, err)
	})

	t.Run("Delete category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)
		repoErr := errors.New("repository delete failed")

		repo.EXPECT().
			DeleteCategory(ctx, 1, 12).
			Return([]core_domain.DeletedTransaction{}, repoErr).
			Times(1)

		err := service.DeleteCategory(ctx, 1, 12)

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("Delete category invalid ID", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := noopBudgetClient{}
		service := NewCategoryService(repo, client)

		err := service.DeleteCategory(ctx, core_domain.UninitializedID, 12)

		assert.Error(t, err)
	})
}
