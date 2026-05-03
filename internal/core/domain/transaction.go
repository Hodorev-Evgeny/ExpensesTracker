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
