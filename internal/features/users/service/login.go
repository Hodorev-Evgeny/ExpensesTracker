package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *UserService) LoginUser(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("invalid user: %w", err)
	}

	userID, err := s.userRepository.FindByEmail(ctx, user)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	if user.ChangeID(userID) != nil {
		return "", fmt.Errorf("invalid user")
	}

	sessionID, err := s.userRedisRepository.CreateCache(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create session cache: %w", err)
	}

	return sessionID, nil
}
