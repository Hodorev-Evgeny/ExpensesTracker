package userservice

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *UserService) GetUserByID(ctx context.Context, id int) (core_domain.User, error) {
	if id == core_domain.UninitializedID {
		return core_domain.User{}, fmt.Errorf("not allowed: %w", core_errors.ErrValidation)
	}

	userDomain, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("error getting user in service: %w", err)
	}

	return userDomain, nil
}
