package grpcclient

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type OptionalArgsStatus struct {
	CategoryID core_domain.Nullable[int64]
	At         core_domain.Nullable[time.Time]
}
