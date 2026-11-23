package categoryservice

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	corelogger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeCategoryBudgetClient struct {
	mu        sync.Mutex
	applied   []core_domain.Transaction
	revoked   []core_domain.DeletedTransaction
	applyErr  error
	revokeErr error
}

func (s *fakeCategoryBudgetClient) ListBudgetStatuses(ctx context.Context, user_id int64, args grpcclient.OptionalArgsStatus) (*pb.ListBudgetStatusesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *fakeCategoryBudgetClient) ApplySpending(
	_ context.Context,
	transaction core_domain.Transaction,
) (*pb.ApplySpendingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applied = append(s.applied, transaction)
	if s.applyErr != nil {
		return nil, s.applyErr
	}
	return &pb.ApplySpendingResponse{}, nil
}

func (s *fakeCategoryBudgetClient) RevokeSpending(
	_ context.Context,
	transaction core_domain.DeletedTransaction,
) (*pb.RevokeSpendingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revoked = append(s.revoked, transaction)
	if s.revokeErr != nil {
		return nil, s.revokeErr
	}
	return &pb.RevokeSpendingResponse{}, nil
}

func (s *fakeCategoryBudgetClient) GetTimeOut() time.Duration {
	return time.Second
}

func (s *fakeCategoryBudgetClient) calls() (
	[]core_domain.Transaction,
	[]core_domain.DeletedTransaction,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]core_domain.Transaction(nil), s.applied...),
		append([]core_domain.DeletedTransaction(nil), s.revoked...)
}

var _ grpcclient.BudgetService = (*fakeCategoryBudgetClient)(nil)

func TestServiceCategory(t *testing.T) {
	ctx := corelogger.ToContext(context.Background(), &corelogger.Logger{
		Logger: zap.NewNop(),
	})

	t.Run("Create correct category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
		service := NewCategoryService(repo, client)

		result, err := service.GetCategoryByID(ctx, core_domain.UninitializedID, 12)

		assert.Error(t, err)
		assert.Equal(t, core_domain.Category{}, result)
	})

	t.Run("Successful rename category", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
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
		client := &fakeCategoryBudgetClient{}
		service := NewCategoryService(repo, client)

		deleted := []core_domain.DeletedTransaction{
			{ID: 101, UserID: 12},
			{ID: 102, UserID: 12},
			{ID: 103, UserID: 12},
		}

		repo.EXPECT().
			DeleteCategory(ctx, 1, 12).
			Return(deleted, nil).
			Times(1)

		err := service.DeleteCategory(ctx, 1, 12)

		require.NoError(t, err)
		applied, revoked := client.calls()
		assert.Empty(t, applied)
		assert.Len(t, revoked, len(deleted))
		assert.ElementsMatch(t, deleted, revoked)
	})

	t.Run("Delete category repository error", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := &fakeCategoryBudgetClient{}
		service := NewCategoryService(repo, client)
		repoErr := errors.New("repository delete failed")

		deleted := []core_domain.DeletedTransaction{
			{ID: 101, UserID: 12},
			{ID: 102, UserID: 12},
		}

		repo.EXPECT().
			DeleteCategory(ctx, 1, 12).
			Return(deleted, repoErr).
			Times(1)

		err := service.DeleteCategory(ctx, 1, 12)

		assert.ErrorIs(t, err, repoErr)
		applied, revoked := client.calls()
		assert.Empty(t, applied)
		assert.Empty(t, revoked)
	})

	t.Run("Delete category Budget unavailable", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := &fakeCategoryBudgetClient{
			revokeErr: errors.New("budget unavailable"),
		}
		service := NewCategoryService(repo, client)

		deleted := []core_domain.DeletedTransaction{
			{ID: 101, UserID: 12},
			{ID: 102, UserID: 12},
		}

		repo.EXPECT().
			DeleteCategory(ctx, 1, 12).
			Return(deleted, nil).
			Times(1)

		err := service.DeleteCategory(ctx, 1, 12)

		assert.NoError(t, err)
		applied, revoked := client.calls()
		assert.Empty(t, applied)
		assert.Len(t, revoked, len(deleted))
		assert.ElementsMatch(t, deleted, revoked)
	})

	t.Run("Delete category invalid ID", func(t *testing.T) {
		repo := NewMockCategoryRepository(t)
		client := &fakeCategoryBudgetClient{}
		service := NewCategoryService(repo, client)

		err := service.DeleteCategory(ctx, core_domain.UninitializedID, 12)

		assert.Error(t, err)
		applied, revoked := client.calls()
		assert.Empty(t, applied)
		assert.Empty(t, revoked)
	})
}
