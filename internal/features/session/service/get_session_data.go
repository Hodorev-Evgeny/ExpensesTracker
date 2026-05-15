package feature_service_session

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *SessionService) GetSessionData(
	ctx context.Context,
	key string,
) (core_domain.CookieData, error) {
	if key == "" {
		return core_domain.CookieData{}, fmt.Errorf("key is empty")
	}

	correctKey := "sessionID:" + key
	cookie, err := s.sessionRepository.GetSessionData(ctx, correctKey)
	if err != nil {
		return core_domain.CookieData{}, fmt.Errorf("error get session data from service: %w", err)
	}

	return cookie, nil
}
