package grpcclient

import (
	"context"
	"fmt"
	"math"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	pb "github.com/diekinari/budget/api/budgetpb/budget/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BudgetService interface {
	ApplySpending(
		ctx context.Context,
		transaction core_domain.Transaction,
	) (*pb.ApplySpendingResponse, error)

	RevokeSpending(
		ctx context.Context,
		transaction core_domain.DeletedTransaction,
	) (*pb.RevokeSpendingResponse, error)
}

type BudgetClient struct {
	client pb.BudgetServiceClient
	conf   BudgetConfig
}

func newBudgetClient(
	conn grpc.ClientConnInterface,
	conf BudgetConfig,
) *BudgetClient {
	return &BudgetClient{
		client: pb.NewBudgetServiceClient(conn),
		conf:   conf,
	}
}

func (c *BudgetClient) ApplySpending(
	ctx context.Context,
	transaction core_domain.Transaction,
) (*pb.ApplySpendingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.conf.Timeout)
	defer cancel()

	amount, err := toMinor(transaction.Sum)
	if err != nil {
		return nil, fmt.Errorf("could not convert transaction amount: %w", err)
	}

	ans, err := c.client.ApplySpending(ctx, &pb.ApplySpendingRequest{
		TransactionId: int64(transaction.ID),
		UserId:        int64(transaction.UserID),
		CategoryId:    int64(transaction.CategoryID),
		Amount: &pb.Money{
			AmountMinor: amount,
			Currency:    "RUB",
		},
		OccurredAt: timestamppb.New(transaction.Date),
	})
	if err != nil {
		return nil, fmt.Errorf("could not apply spending request: %w", err)
	}

	return ans, nil
}

func (c *BudgetClient) RevokeSpending(
	ctx context.Context,
	transaction core_domain.DeletedTransaction,
) (*pb.RevokeSpendingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.conf.Timeout)
	defer cancel()

	ans, err := c.client.RevokeSpending(ctx, &pb.RevokeSpendingRequest{
		TransactionId: int64(transaction.ID),
		UserId:        int64(transaction.UserID),
	})
	if err != nil {
		return nil, fmt.Errorf("could not revoke spending transaction: %w", err)
	}

	return ans, nil
}

func NewClientBudget(conf BudgetConfig) (*BudgetClient, error) {
	conn, err := grpc.NewClient(
		conf.ADDR,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return newBudgetClient(conn, conf), nil
}

func toMinor(sum int) (int64, error) {
	amount := int64(sum)

	if amount > math.MaxInt64/100 ||
		amount < math.MinInt64/100 {
		return 0, fmt.Errorf("amount overflow: %d", amount)
	}

	return amount * 100, nil
}
