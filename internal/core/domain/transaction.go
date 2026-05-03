package core_domain

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID          int
	Sum         int
	Type        string
	Date        time.Time
	CategoryID  int
	UserID      int
	Comments    string
	TimeCreated time.Time
	TimeChange  *time.Time
}

type FiltersTransaction struct {
	Limit      *int       `json:"limit"`
	Offset     *int       `json:"offset"`
	UserId     *int       `json:"user_id"`
	CategoryId *int       `json:"category_id"`
	Sum        *int       `json:"sum"`
	From       *time.Time `json:"from"`
	To         *time.Time `json:"to"`
}

func NewTransactionUnincelized(
	Sum, CategoryID, UserID int,
	Type, Comments string,
	Date, TimeCreated time.Time,
	TimeChange *time.Time,
) Transaction {
	return Transaction{
		ID:          UnincelizedID,
		Sum:         Sum,
		CategoryID:  CategoryID,
		UserID:      UserID,
		Type:        Type,
		Date:        Date,
		TimeCreated: TimeCreated,
		TimeChange:  TimeChange,
	}
}

func (f *FiltersTransaction) Validate() error {
	if (f.Limit != nil && *f.Limit < 0) || (f.Offset != nil && *f.Offset < 0) {
		return fmt.Errorf("invalid filter parameters limit/offset")
	}

	if f.Sum != nil && *f.Sum < 0 {
		return fmt.Errorf("invalid filter parameters sum")
	}

	if (f.UserId != nil && *f.UserId == UnincelizedID) || (f.CategoryId != nil && *f.CategoryId == UnincelizedID) {
		return fmt.Errorf("invalid filter parameters categoryID/userID")
	}

	return nil
}

func (t *Transaction) Validate() error {
	if t.Sum <= 0 {
		return fmt.Errorf("invalid transaction sum: %d", t.Sum)
	}

	if t.CategoryID == UnincelizedID || t.UserID == UnincelizedID {
		return fmt.Errorf("invalid transaction category id: %d", t.CategoryID)
	}

	if len(t.Comments) >= 1000 {
		return fmt.Errorf("invalid transaction comments: %d", len(t.Comments))
	}

	return nil
}
