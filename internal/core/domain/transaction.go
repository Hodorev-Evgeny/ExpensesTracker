package core_domain

import (
	"fmt"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
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

type PatchTransaction struct {
	Sum             Nullable[int]
	CategoryId      Nullable[int]
	TypeTransaction Nullable[string]
	Comments        Nullable[string]
}

func NewTransactionPatch(
	sum Nullable[int],
	category Nullable[int],
	typeTransaction Nullable[string],
	comments Nullable[string],
) PatchTransaction {
	return PatchTransaction{
		Sum:             sum,
		CategoryId:      category,
		TypeTransaction: typeTransaction,
		Comments:        comments,
	}
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
		Comments:    Comments,
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

	if t.Type != "Income" && t.Type != "Expenditure" {
		return fmt.Errorf("invalid transaction type: %s", t.Type)
	}

	if len(t.Comments) >= 1000 {
		return fmt.Errorf("invalid transaction comments: %d", len(t.Comments))
	}

	return nil
}

func (p *PatchTransaction) Validate() error {
	if p.Sum.Set && *p.Sum.Value <= 0 {
		return fmt.Errorf("invalid transaction sum: %d", p.Sum.Value)
	}

	if p.CategoryId.Set && *p.CategoryId.Value == UnincelizedID {
		return fmt.Errorf("invalid transaction category id: %d", p.CategoryId.Value)
	}

	if p.TypeTransaction.Set && *p.TypeTransaction.Value != "Income" && *p.TypeTransaction.Value != "Expenditure" {
		return fmt.Errorf("invalid transaction type: %s", p.TypeTransaction.Value)
	}

	if p.Comments.Set && len(*p.Comments.Value) >= 1000 {
		return fmt.Errorf("invalid transaction comments: %d", p.Comments.Value)
	}

	return nil
}

func (t *Transaction) ApplyPatch(patch PatchTransaction) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", core_errors.ErrorValidation)
	}

	tmp := *t

	if patch.TypeTransaction.Set {
		tmp.Type = *patch.TypeTransaction.Value
	}

	if patch.CategoryId.Set {
		tmp.CategoryID = *patch.CategoryId.Value
	}

	if patch.Sum.Set {
		tmp.Sum = *patch.Sum.Value
	}

	if patch.Comments.Set {
		tmp.Comments = *patch.Comments.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}

	timeChanged := time.Now()
	tmp.TimeChange = &timeChanged

	*t = tmp

	return nil
}
