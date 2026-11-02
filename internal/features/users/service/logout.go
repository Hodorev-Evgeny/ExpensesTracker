package feature_user_service

import (
	"context"
	"fmt"
)

func (s *UserService) LogoutUser(
	ctx context.Context,
	sessionID string,
) error {
	err := s.userRedisRepository.DeleteCache(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("logout user: %w", err)
	}

	return nil
}
