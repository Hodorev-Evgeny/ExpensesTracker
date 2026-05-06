package feature_transport_limit

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type limitDTO struct {
	Duration    time.Time `json:"duration" validate:"required"`
	AmountLimit int       `json:"amount_limit" validate:"required"`
}

type LimitResponse struct {
	ID int `json:"id"`
	limitDTO
}

func LimitResponseToDomain(req limitDTO) core_domain.Limit {
	return core_domain.CreateUnincelizeLimit(
		req.Duration,
		req.AmountLimit,
	)
}

func LimitDomainToResponse(limit core_domain.Limit) LimitResponse {
	return LimitResponse{
		ID: limit.ID,
		limitDTO: limitDTO{
			Duration:    limit.Duration,
			AmountLimit: limit.AmountLimit,
		},
	}
}
