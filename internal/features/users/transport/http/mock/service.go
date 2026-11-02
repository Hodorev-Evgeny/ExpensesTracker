package transport_user_mock

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type MockService struct {
}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) CreateUser(
	ctx context.Context,
	req core_domain.User,
) (core_domain.User, error) {
	return core_domain.User{}, nil
}

func (s *MockService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]core_domain.User, error) {
	return []core_domain.User{}, nil
}

func (s *MockService) ExtraditionUser(
	ctx context.Context,
	id int,
) (core_domain.User, error) {
	return core_domain.User{}, nil
}

func (s *MockService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	return nil
}

func (s *MockService) PatchUser(
	ctx context.Context,
	id int,
	req core_domain.UserPatch,
) (core_domain.User, error) {
	return core_domain.User{}, nil
}

func (s *MockService) CreateCache(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	return "", nil
}

func (s *MockService) LoginUser(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	return "", nil
}

func (s *MockService) LogoutUser(
	ctx context.Context,
	sessionID string,
) error {
	return nil
}
