package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) CreateUser(ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {

	if err := user.Validate(); err != nil {
		return core_domain.User{}, fmt.Errorf("%w: %w", core_errors.ErrorValidation, err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hash)

	var initializedUser core_domain.User
	initializedUser, err = s.userRepository.AddUser(ctx, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("error when try to create user: %w", err)
	}

	return initializedUser, nil
}
