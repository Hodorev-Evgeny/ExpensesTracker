package transactionservice

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceTransaction(t *testing.T) {
	ctx := context.Background()

	newTransaction := func() core_domain.Transaction {
		return core_domain.Transaction{
			ID:         1,
			Sum:        1000,
			Type:       "Expenditure",
			CategoryID: 2,
			UserID:     12,
			Comments:   "Test payment",
		}
	}

	t.Run("Create correct transaction", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()

		repo.EXPECT().
			CheckCategoryOwner(ctx, transaction.CategoryID, transaction.UserID).
			Return(true, nil).
			Times(1)

		repo.EXPECT().
			CreateTransaction(ctx, transaction).
			Return(transaction, nil).
			Times(1)

		result, err := service.CreateTransaction(ctx, transaction)

		assert.NoError(t, err)
		assert.Equal(t, transaction, result)
	})

	t.Run("Create invalid transaction", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()
		transaction.Sum = -100

		result, err := service.CreateTransaction(ctx, transaction)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Create transaction repository error", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()
		repoErr := core_repository_pool.ErrViolationKey

		repo.EXPECT().
			CheckCategoryOwner(ctx, transaction.CategoryID, transaction.UserID).
			Return(true, nil).
			Times(1)

		repo.EXPECT().
			CreateTransaction(ctx, transaction).
			Return(core_domain.Transaction{}, repoErr).
			Times(1)

		result, err := service.CreateTransaction(ctx, transaction)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Get existing transaction", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()

		repo.EXPECT().
			GetTransaction(ctx, transaction.ID, transaction.UserID).
			Return(transaction, nil).
			Times(1)

		result, err := service.GetTransaction(ctx, transaction.ID, transaction.UserID)

		assert.NoError(t, err)
		assert.Equal(t, transaction, result)
	})

	t.Run("Get transaction repository error", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		repoErr := errors.New("transaction not found")

		repo.EXPECT().
			GetTransaction(ctx, 1, 12).
			Return(core_domain.Transaction{}, repoErr).
			Times(1)

		result, err := service.GetTransaction(ctx, 1, 12)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Get transaction invalid ID", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)

		result, err := service.GetTransaction(ctx, core_domain.UnincelizedID, 12)

		assert.Error(t, err)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Get transactions", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		limit, offset := 10, 0
		filters := core_domain.FiltersTransaction{Limit: &limit, Offset: &offset}
		expected := []core_domain.Transaction{newTransaction()}

		repo.EXPECT().
			GetTransactions(ctx, filters).
			Return(expected, nil).
			Times(1)

		result, err := service.GetTransactions(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Get transactions repository error", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		limit := 10
		filters := core_domain.FiltersTransaction{Limit: &limit}
		repoErr := errors.New("repository list failed")

		repo.EXPECT().
			GetTransactions(ctx, filters).
			Return(nil, repoErr).
			Times(1)

		result, err := service.GetTransactions(ctx, filters)

		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("Get transactions invalid filters", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		limit := -1
		filters := core_domain.FiltersTransaction{Limit: &limit}

		result, err := service.GetTransactions(ctx, filters)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
		assert.Nil(t, result)
	})

	t.Run("Successful patch transaction", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()
		comments := "Updated payment"
		patch := core_domain.PatchTransaction{
			Comments: core_domain.Nullable[string]{Value: &comments, Set: true},
		}
		expected := transaction
		expected.Comments = comments

		repo.EXPECT().
			GetTransaction(ctx, transaction.ID, transaction.UserID).
			Return(transaction, nil).
			Times(1)

		repo.EXPECT().
			PatchTransaction(ctx, mock.MatchedBy(func(got core_domain.Transaction) bool {
				return got.ID == expected.ID &&
					got.Comments == expected.Comments &&
					got.Sum == expected.Sum &&
					got.TimeChange != nil
			}), transaction.UserID).
			Return(expected, nil).
			Times(1)

		result, err := service.PatchTransaction(ctx, transaction.ID, patch, transaction.UserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Patch transaction get error", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		repoErr := errors.New("transaction not found")
		comments := "Updated payment"
		patch := core_domain.PatchTransaction{
			Comments: core_domain.Nullable[string]{Value: &comments, Set: true},
		}

		repo.EXPECT().
			GetTransaction(ctx, 1, 12).
			Return(core_domain.Transaction{}, repoErr).
			Times(1)

		result, err := service.PatchTransaction(ctx, 1, patch, 12)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Patch transaction invalid data", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()
		invalidType := "unknown"
		patch := core_domain.PatchTransaction{
			TypeTransaction: core_domain.Nullable[string]{Value: &invalidType, Set: true},
		}

		repo.EXPECT().
			GetTransaction(ctx, transaction.ID, transaction.UserID).
			Return(transaction, nil).
			Times(1)

		result, err := service.PatchTransaction(ctx, transaction.ID, patch, transaction.UserID)

		assert.ErrorIs(t, err, core_errors.ErrorValidation)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Patch transaction repository error", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		transaction := newTransaction()
		comments := "Updated payment"
		patch := core_domain.PatchTransaction{
			Comments: core_domain.Nullable[string]{Value: &comments, Set: true},
		}
		repoErr := errors.New("repository patch failed")

		repo.EXPECT().
			GetTransaction(ctx, transaction.ID, transaction.UserID).
			Return(transaction, nil).
			Times(1)

		repo.EXPECT().
			PatchTransaction(ctx, mock.MatchedBy(func(got core_domain.Transaction) bool {
				return got.ID == transaction.ID &&
					got.Comments == comments &&
					got.TimeChange != nil
			}), transaction.UserID).
			Return(core_domain.Transaction{}, repoErr).
			Times(1)

		result, err := service.PatchTransaction(ctx, transaction.ID, patch, transaction.UserID)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.Transaction{}, result)
	})

	t.Run("Successful delete transaction", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)

		repo.EXPECT().
			DeleteTransaction(ctx, 1, 12).
			Return(nil).
			Times(1)

		err := service.DeleteTransaction(ctx, 1, 12)

		assert.NoError(t, err)
	})

	t.Run("Delete transaction repository error", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)
		repoErr := errors.New("repository delete failed")

		repo.EXPECT().
			DeleteTransaction(ctx, 1, 12).
			Return(repoErr).
			Times(1)

		err := service.DeleteTransaction(ctx, 1, 12)

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("Delete transaction invalid ID", func(t *testing.T) {
		repo := NewMockTransactionRepository(t)
		service := NewTransactionService(repo)

		err := service.DeleteTransaction(ctx, core_domain.UnincelizedID, 12)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
	})
}
