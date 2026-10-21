package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) LoginUser(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("invalid user: %w", err)
	}

	userLogin, err := s.userRepository.FindByEmail(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error get user for login: %w", err)
	}

	errHash := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(user.Password))
	if errHash != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	sessionID, err := s.userRedisRepository.CreateCache(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create session cache: %w", err)
	}

	return sessionID, nil
}
