package feature_transactio_transport

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type TransactionDTO struct {
	Sum        int        `json:"sum" validate:"required"`
	Type       string     `json:"type" validate:"required"`
	Date       time.Time  `json:"date" validate:"required"`
	UserID     int        `json:"user_id" validate:"required"`
	CategoryID int        `json:"category_id" validate:"required"`
	Comments   string     `json:"comments" validate:"required"`
	TimeChange *time.Time `json:"time_change"`
}

type TransactionResponse struct {
	ID         int        `json:"id"`
	Sum        int        `json:"sum" validate:"required"`
	Type       string     `json:"type" validate:"required"`
	Date       time.Time  `json:"date" validate:"required"`
	UserID     int        `json:"user_id" validate:"required"`
	CategoryID int        `json:"category_id" validate:"required"`
	Comments   string     `json:"comments" validate:"required"`
	TimeChange *time.Time `json:"time_change"`
}

func ToDomainTransaction(
	req TransactionRequest,
) core_domain.Transaction {
	return core_domain.NewTransactionUnincelized(
		req.Sum, req.CategoryID, req.UserID,
		req.Type, req.Comments,
		req.Date, time.Now(),
		req.TimeChange,
	)
}

func ToTransactionResponse(
	transaction core_domain.Transaction,
) TransactionResponse {
	return TransactionResponse{
		ID:         transaction.ID,
		Sum:        transaction.Sum,
		Type:       transaction.Type,
		Date:       transaction.Date,
		UserID:     transaction.UserID,
		CategoryID: transaction.CategoryID,
		Comments:   transaction.Comments,
		TimeChange: transaction.TimeChange,
	}
}

func ToDomainPatch(patch TransactionPatchRequest) core_domain.PatchTransaction {
	return core_domain.NewTransactionPatch(
		patch.Sum.ToDomain(),
		patch.CategoryID.ToDomain(),
		patch.TypeTransaction.ToDomain(),
		patch.Comments.ToDomain(),
	)
}
