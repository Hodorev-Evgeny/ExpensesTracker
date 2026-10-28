package feature_user_service

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) LoginUser(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	if err := user.LoginValidate(); err != nil {
		return "", fmt.Errorf("invalid user: %w", err)
	}

	userLogin, err := s.userRepository.FindByEmail(ctx, user)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return "", fmt.Errorf("error get user for login: %w", core_errors.ErrorUnauthorized)
		}
		return "", fmt.Errorf("error find user by email: %w", err)
	}

	errHash := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(user.Password))
	if errHash != nil {
		return "", fmt.Errorf("the password does not match: %w", core_errors.ErrorUnauthorized)
	}

	sessionID, err := s.userRedisRepository.CreateCache(ctx, userLogin)
	if err != nil {
		return "", fmt.Errorf("create session cache: %w", err)
	}

	return sessionID, nil
}
