package core_domain

import (
	"fmt"
	"time"
)

type Limit struct {
	ID          int
	Duration    time.Time
	AmountLimit int
}

func NewLimit(
	duration time.Time,
	amountLimit int,
) Limit {
	return Limit{
		ID:          UnincelizedID,
		Duration:    duration,
		AmountLimit: amountLimit,
	}
}

func CreateUnincelizeLimit(
	duration time.Time,
	amountLimit int,
) Limit {
	return NewLimit(duration, amountLimit)
}

func (l *Limit) Validate() error {
	if l.AmountLimit <= 0 {
		return fmt.Errorf("amount limit must be greater than zero")
	}

	return nil
}
