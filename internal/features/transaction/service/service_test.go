package feature_transaction_service

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

type MockRepository struct{}

func (r *MockRepository) CreateTransaction(
	ctx context.Context,
	transaction core_domain.Transaction,
) (core_domain.Transaction, error) {
	return core_domain.Transaction{}, nil
}

func (r *MockRepository) GetTransactions(
	ctx context.Context,
	filters core_domain.FiltersTransaction,
) ([]core_domain.Transaction, error) {
	return []core_domain.Transaction{}, nil
}

func (r *MockRepository) GetTransaction(
	ctx context.Context,
	id int,
	userID int,
) (core_domain.Transaction, error) {
	return core_domain.Transaction{}, nil
}

func (r *MockRepository) PatchTransaction(
	ctx context.Context,
	transaction core_domain.Transaction,
	userID int,
) (core_domain.Transaction, error) {
	return core_domain.Transaction{}, nil
}

func (r *MockRepository) DeleteTransaction(
	ctx context.Context,
	id int,
	userID int,
) error {
	return nil
}

func (r *MockRepository) CheckCategoryOwner(
	ctx context.Context,
	categoryID int,
	userID int,
) (bool, error) {
	return false, nil
}

func TestCheckCategoryOwner(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	seevice := NewTransactionService(repo)

	t.Run("Test create trunsaction", func(t *testing.T) {
		transaction := core_domain.Transaction{
			ID:          1,
			Sum:         12,
			Type:        "Income",
			Date:        time.Now(),
			CategoryID:  2,
			UserID:      1,
			Comments:    "",
			TimeCreated: time.Now(),
		}

		_, err := seevice.CreateTransaction(ctx, transaction)
		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Errorf("expected error not foud, got: %v", err)
		}
	})

	t.Run("Test patch transaction", func(t *testing.T) {
		transaction := core_domain.PatchTransaction{
			CategoryId: struct {
				Value *int
				Set   bool
			}{Value: new(2), Set: true},
		}

		_, err := seevice.PatchTransaction(ctx, 1, transaction, 1)
		if !errors.Is(err, core_errors.ErrorNotFoud) {
			t.Errorf("expected error not foud, got: %v", err)
		}
	})
}
