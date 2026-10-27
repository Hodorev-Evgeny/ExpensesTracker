package feature_user_service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"golang.org/x/crypto/bcrypt"
)

var testUserInfo map[string]*core_domain.User = map[string]*core_domain.User{}

type TestCase struct {
	name            string
	inputUser       core_domain.User
	expectedError   error
	withCreateCache bool
}

type RepositoryRedisTestMock struct {
	IncomingUser core_domain.User
}

func (r *RepositoryRedisTestMock) CreateCache(
	ctx context.Context,
	user core_domain.User,
) (string, error) {
	r.IncomingUser = user
	return "", nil
}

func (r *RepositoryRedisTestMock) DeleteCache(
	ctx context.Context,
	sessionID string,
) error {
	return nil
}

type RepositoryPostgresMock struct{}

func (r *RepositoryPostgresMock) AddUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	return core_domain.User{}, nil
}

func (r *RepositoryPostgresMock) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]core_domain.User, error) {
	return []core_domain.User{}, nil
}

func (r *RepositoryPostgresMock) ExtraditionUser(
	ctx context.Context,
	id int,
) (core_domain.User, error) {
	return core_domain.User{}, nil
}

func (r *RepositoryPostgresMock) DeleteUser(
	ctx context.Context,
	id int,
) error {
	return nil
}

func (r *RepositoryPostgresMock) PatchUser(
	ctx context.Context,
	id int,
	patch core_domain.User,
) (core_domain.User, error) {
	return core_domain.User{}, nil
}

func (r *RepositoryPostgresMock) FindByEmail(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	if logUser, ok := testUserInfo[user.Email]; ok {
		return *logUser, nil
	}
	return core_domain.User{}, fmt.Errorf("user not found: %w", core_repository_pool.ErrNoRows)
}

func HashTestPassword() error {
	for _, user := range testUserInfo {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash test password: %w", err)
		}
		user.Password = string(hash)
	}
	return nil
}

func TestUserLogin(t *testing.T) {
	testCase := []TestCase{
		{
			name: "Correct Login",
			inputUser: core_domain.User{
				ID:       -1,
				Email:    "testcorrect@gmai.com",
				Password: "VeryHardAndCorrectPass",
			},
			expectedError:   nil,
			withCreateCache: true,
		}, {
			name: "Incorrect password",
			inputUser: core_domain.User{
				ID:       -1,
				Email:    "testcorrect@gmai.com",
				Password: "VeryHardAndPass",
			},
			expectedError:   core_errors.ErrorUnauthorized,
			withCreateCache: false,
		},
		{
			name: "Unknown email",
			inputUser: core_domain.User{
				ID:       -1,
				Email:    "testunknow@gmai.com",
				Password: "VeryHardAndCorrectPass",
			},
			expectedError:   core_errors.ErrorUnauthorized,
			withCreateCache: false,
		},
	}
	testUserInfo["testcorrect@gmai.com"] = &core_domain.User{
		ID:       123,
		Email:    "testcorrect@gmai.com",
		Password: "VeryHardAndCorrectPass",
	}

	if err := HashTestPassword(); err != nil {
		t.Error(err)
	}

	for _, test := range testCase {
		redisMock := &RepositoryRedisTestMock{}
		postgresMock := &RepositoryPostgresMock{}
		service := NewUserService(postgresMock, redisMock)

		_, err := service.LoginUser(context.Background(), test.inputUser)

		if !errors.Is(err, test.expectedError) {
			t.Errorf("Test case %s expected error: %s, got: %s", test.name, test.expectedError, err)
		}

		if expectedUser, ok := testUserInfo[test.inputUser.Email]; test.withCreateCache && ok {
			if expectedUser.ID != redisMock.IncomingUser.ID {
				t.Errorf("Test case %s expected user ID: %d, got: %d", test.name, expectedUser.ID, redisMock.IncomingUser.ID)
			}
		}
	}
}
