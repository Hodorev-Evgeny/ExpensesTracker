package feature_service_session

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type SessionRepository interface {
	GetSessionData(
		ctx context.Context,
		key string,
	) (core_domain.CookieData, error)
}

type SessionService struct {
	sessionRepository SessionRepository
}

func NewSessionService(
	sessionRepository SessionRepository,
) *SessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
	}
}
