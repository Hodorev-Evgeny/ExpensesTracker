package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *UserService) CreateCache(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	sessionID, err := s.userRedisRepository.CreateCache(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create session id: %w", err)
	}

	return sessionID, nil
}
