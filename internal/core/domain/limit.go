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
}

type PatchLimit struct {
	Duraction    Nullable[time.Time]
	Amount_limit Nullable[int]
}

func NewLimit(
	id int,
	duration time.Time,
	amountLimit int,
) Limit {
	return Limit{
		ID:          id,
		Duration:    duration,
		AmountLimit: amountLimit,
	}
}

func CreateUnincelizeLimit(
	duration time.Time,
	amountLimit int,
) Limit {
	return NewLimit(UnincelizedID, duration, amountLimit)
}

func CreatePatchLimit(
	duration Nullable[time.Time],
	amountLimit Nullable[int],
) PatchLimit {
	return PatchLimit{
		Duraction:    duration,
		Amount_limit: amountLimit,
	}
}

func (p *PatchLimit) Validate() error {
	if p.Amount_limit.Set && *p.Amount_limit.Value < 0 {
		return fmt.Errorf("amount limit cannot be set to zero")
	}

	return nil
}

func (l *Limit) Validate() error {
	if l.AmountLimit <= 0 {
		return fmt.Errorf("amount limit must be greater than zero")
	}

	return nil
}

func (l *Limit) ApplyPatch(patch PatchLimit) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", core_errors.ErrorValidation)
	}

	tmp := *l
	if patch.Duraction.Set {
		tmp.Duration = *patch.Duraction.Value
	}

	if patch.Amount_limit.Set {
		tmp.AmountLimit = *patch.Amount_limit.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("invalid limit: %w", core_errors.ErrorValidation)
	}

	*l = tmp

	return nil
}
