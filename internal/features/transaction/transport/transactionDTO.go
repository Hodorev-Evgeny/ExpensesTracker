package feature_transactio_transport

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type TransactionDTO struct {
	Sum        int       `json:"sum" validate:"required" example:"10123"`
	Type       string    `json:"type" validate:"required" example:"Expenditure/Income"`
	Date       time.Time `json:"date" validate:"required" example:"2006-01-02T15:04:05-07:00"`
	UserID     int       `json:"user_id" validate:"required" example:"1"`
	CategoryID int       `json:"category_id" validate:"required" example:"1"`
	Comments   string    `json:"comments" validate:"required" example:"somthing else"`
}

type TransactionResponse struct {
	ID         int        `json:"id" example:"1"`
	Sum        int        `json:"sum" validate:"required" example:"10123"`
	Type       string     `json:"type" validate:"required" example:"Expenditure/Income"`
	Date       time.Time  `json:"date" validate:"required" example:"2006-01-02T15:04:05-07:00"`
	UserID     int        `json:"user_id" validate:"required" example:"1"`
	CategoryID int        `json:"category_id" validate:"required" example:"1"`
	Comments   string     `json:"comments" validate:"required" example:"somthing else"`
	TimeChange *time.Time `json:"time_change" example:"2006-01-02T15:04:05-07:00"`
}

func ToDomainTransaction(
	req TransactionRequest,
) core_domain.Transaction {
	return core_domain.NewTransactionUnincelized(
		req.Sum, req.CategoryID, req.UserID,
		req.Type, req.Comments,
		req.Date, time.Now(),
		nil,
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
