package core_domain

import (
	"fmt"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

type Limit struct {
	ID          int
	Duration    time.Time
	AmountLimit int
	UserID      int
}

type PatchLimit struct {
	Duration     Nullable[time.Time]
	Amount_limit Nullable[int]
}

func NewLimit(
	id int,
	duration time.Time,
	amountLimit int,
	userID int,
) Limit {
	return Limit{
		ID:          id,
		Duration:    duration,
		AmountLimit: amountLimit,
		UserID:      userID,
	}
}

func CreateUninitializedLimit(
	duration time.Time,
	amountLimit int,
	userID int,
) Limit {
	return NewLimit(UninitializedID, duration, amountLimit, userID)
}

func CreatePatchLimit(
	duration Nullable[time.Time],
	amountLimit Nullable[int],
) PatchLimit {
	return PatchLimit{
		Duration:     duration,
		Amount_limit: amountLimit,
	}
}

func (p *PatchLimit) Validate() error {
	if p.Amount_limit.Set && *p.Amount_limit.Value < 0 {
		return fmt.Errorf("amount limit cannot be set to zero: %w", core_errors.ErrBadRequest)
	}

	return nil
}

func (l *Limit) Validate() error {
	if l.AmountLimit <= 0 {
		return fmt.Errorf("amount limit must be greater than zero: %w", core_errors.ErrBadRequest)
	}

	if l.Duration.IsZero() {
		return fmt.Errorf("invalid duration: %w", core_errors.ErrBadRequest)
	}

	return nil
}

func (l *Limit) ApplyPatch(patch PatchLimit) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", core_errors.ErrValidation)
	}

	tmp := *l
	if patch.Duration.Set {
		tmp.Duration = *patch.Duration.Value
	}

	if patch.Amount_limit.Set {
		tmp.AmountLimit = *patch.Amount_limit.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("invalid limit: %w", core_errors.ErrValidation)
	}

	*l = tmp

	return nil
}
