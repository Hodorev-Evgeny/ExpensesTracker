package feature_transport_limit

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type LimitDTO struct {
	Duration    time.Time `json:"duration" validate:"required" example:"2006-01-02T15:04:05-07:00"`
	AmountLimit int       `json:"amount_limit" validate:"required" example:"1000"`
}

type LimitResponse struct {
	ID int `json:"id"`
	LimitDTO
}

func LimitResponseToDomain(req LimitDTO) core_domain.Limit {
	return core_domain.CreateUnincelizeLimit(
		req.Duration,
		req.AmountLimit,
	)
}

func LimitResponseToPatch(req PatchLimit) core_domain.PatchLimit {
	return core_domain.CreatePatchLimit(req.Duraction.ToDomain(), req.AmountLimit.ToDomain())
}

func LimitDomainToResponse(limit core_domain.Limit) LimitResponse {
	return LimitResponse{
		ID: limit.ID,
		LimitDTO: LimitDTO{
			Duration:    limit.Duration,
			AmountLimit: limit.AmountLimit,
		},
	}
}
