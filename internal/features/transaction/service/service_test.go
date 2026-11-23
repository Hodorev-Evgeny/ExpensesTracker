package transactionservice

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
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

type FakeBudgetClient struct {
	mu        sync.Mutex
	applied   []core_domain.Transaction
	revoked   []core_domain.DeletedTransaction
	applyErr  error
	revokeErr error
}

func (s *FakeBudgetClient) ListBudgetStatuses(
	ctx context.Context,
	user_id int64,
	args grpcclient.OptionalArgsStatus,
) (*pb.ListBudgetStatusesResponse, error) {
	return &pb.ListBudgetStatusesResponse{}, nil
}

func (s *FakeBudgetClient) GetTimeOut() time.Duration {
	return time.Second
}

func (s *FakeBudgetClient) ApplySpending(
	_ context.Context,
	tx core_domain.Transaction,
) (*pb.ApplySpendingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.applied = append(s.applied, tx)

	if s.applyErr != nil {
		return nil, s.applyErr
	}

	return &pb.ApplySpendingResponse{}, nil
}

func (s *FakeBudgetClient) RevokeSpending(
	_ context.Context,
	tx core_domain.DeletedTransaction,
) (*pb.RevokeSpendingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.revoked = append(s.revoked, tx)

	if s.revokeErr != nil {
		return nil, s.revokeErr
	}

	return &pb.RevokeSpendingResponse{}, nil
}

func TestCheckCategoryOwner(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	seevice := NewTransactionService(repo, &FakeBudgetClient{})

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
		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Errorf("expected error not found, got: %v", err)
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
		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Errorf("expected error not found, got: %v", err)
		}
	})
}
