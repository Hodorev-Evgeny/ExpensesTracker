package feature_user_service

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestServiceUser(t *testing.T) {
	ctx := context.Background()

	newUser := func() core_domain.User {
		return core_domain.User{
			ID:        core_domain.UnincelizedID,
			Full_name: "Test User",
			Email:     "test@example.com",
			Password:  "password123",
			Time_add:  time.Now().UTC(),
		}
	}

	t.Run("Create correct user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		expected := user
		expected.ID = 1

		repo.EXPECT().
			AddUser(ctx, mock.MatchedBy(func(got core_domain.User) bool {
				return got.ID == user.ID &&
					got.Full_name == user.Full_name &&
					got.Email == user.Email &&
					bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(user.Password)) == nil
			})).
			RunAndReturn(func(_ context.Context, got core_domain.User) (core_domain.User, error) {
				got.ID = 1
				return got, nil
			}).
			Times(1)

		result, err := service.CreateUser(ctx, user)

		assert.NoError(t, err)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)))
		expected.Password = result.Password
		assert.Equal(t, expected, result)
	})

	t.Run("Create invalid user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.Full_name = ""

		result, err := service.CreateUser(ctx, user)

		assert.ErrorIs(t, err, core_errors.ErrorValidation)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Create user repository error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		repoErr := errors.New("repository create failed")

		repo.EXPECT().
			AddUser(ctx, mock.MatchedBy(func(got core_domain.User) bool {
				return got.Email == user.Email &&
					bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(user.Password)) == nil
			})).
			Return(core_domain.User{}, repoErr).
			Times(1)

		result, err := service.CreateUser(ctx, user)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Get users", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		limit, offset := 10, 0
		user := newUser()
		user.ID = 1
		expected := []core_domain.User{user}

		repo.EXPECT().
			GetUsers(ctx, &limit, &offset).
			Return(expected, nil).
			Times(1)

		result, err := service.GetUsers(ctx, &limit, &offset)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Get users repository error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		repoErr := errors.New("repository get users failed")

		repo.EXPECT().
			GetUsers(ctx, (*int)(nil), (*int)(nil)).
			Return(nil, repoErr).
			Times(1)

		result, err := service.GetUsers(ctx, nil, nil)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, []core_domain.User{}, result)
	})

	t.Run("Get users invalid limit", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		limit := -1

		result, err := service.GetUsers(ctx, &limit, nil)

		assert.ErrorIs(t, err, core_errors.ErrorValidation)
		assert.Equal(t, []core_domain.User{}, result)
	})

	t.Run("Get users invalid offset", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		offset := -1

		result, err := service.GetUsers(ctx, nil, &offset)

		assert.ErrorIs(t, err, core_errors.ErrorValidation)
		assert.Equal(t, []core_domain.User{}, result)
	})

	t.Run("Get existing user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.ID = 1

		repo.EXPECT().
			ExtraditionUser(ctx, user.ID).
			Return(user, nil).
			Times(1)

		result, err := service.ExtraditionUser(ctx, user.ID)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Get user repository error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		repoErr := errors.New("user not found")

		repo.EXPECT().
			ExtraditionUser(ctx, 1).
			Return(core_domain.User{}, repoErr).
			Times(1)

		result, err := service.ExtraditionUser(ctx, 1)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Get user invalid ID", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)

		result, err := service.ExtraditionUser(ctx, core_domain.UnincelizedID)

		assert.ErrorIs(t, err, core_errors.ErrorValidation)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Successful patch user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.ID = 1
		name := "Updated User"
		patch := core_domain.UserPatch{
			Full_name: core_domain.Nullable[string]{Value: &name, Set: true},
		}
		expected := user
		expected.Full_name = name

		repo.EXPECT().
			ExtraditionUser(ctx, user.ID).
			Return(user, nil).
			Times(1)

		repo.EXPECT().
			PatchUser(ctx, user.ID, expected).
			Return(expected, nil).
			Times(1)

		result, err := service.PatchUser(ctx, user.ID, patch)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Patch user get error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		repoErr := errors.New("user not found")
		name := "Updated User"
		patch := core_domain.UserPatch{
			Full_name: core_domain.Nullable[string]{Value: &name, Set: true},
		}

		repo.EXPECT().
			ExtraditionUser(ctx, 1).
			Return(core_domain.User{}, repoErr).
			Times(1)

		result, err := service.PatchUser(ctx, 1, patch)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Patch user invalid data", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.ID = 1
		patch := core_domain.UserPatch{
			Full_name: core_domain.Nullable[string]{Set: true},
		}

		repo.EXPECT().
			ExtraditionUser(ctx, user.ID).
			Return(user, nil).
			Times(1)

		result, err := service.PatchUser(ctx, user.ID, patch)

		assert.Error(t, err)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Patch user repository error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.ID = 1
		name := "Updated User"
		patch := core_domain.UserPatch{
			Full_name: core_domain.Nullable[string]{Value: &name, Set: true},
		}
		expected := user
		expected.Full_name = name
		repoErr := errors.New("repository patch failed")

		repo.EXPECT().
			ExtraditionUser(ctx, user.ID).
			Return(user, nil).
			Times(1)
		repo.EXPECT().
			PatchUser(ctx, user.ID, expected).
			Return(core_domain.User{}, repoErr).
			Times(1)

		result, err := service.PatchUser(ctx, user.ID, patch)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, core_domain.User{}, result)
	})

	t.Run("Successful delete user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)

		repo.EXPECT().
			DeleteUser(ctx, 1).
			Return(nil).
			Times(1)

		err := service.DeleteUser(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("Delete user repository error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		repoErr := errors.New("repository delete failed")

		repo.EXPECT().
			DeleteUser(ctx, 1).
			Return(repoErr).
			Times(1)

		err := service.DeleteUser(ctx, 1)

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("Delete user invalid ID", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)

		err := service.DeleteUser(ctx, core_domain.UnincelizedID)

		assert.Error(t, err)
	})

	t.Run("Successful login user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.ID = 0

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			t.Fatal(err)
		}
		cachedUser := user
		cachedUser.ID = 12
		cachedUser.Password = string(hash)

		repo.EXPECT().
			FindByEmail(ctx, user).
			Return(cachedUser, nil).
			Times(1)
		redis.EXPECT().
			CreateCache(ctx, cachedUser).
			Return("session-123", nil).
			Times(1)

		sessionID, err := service.LoginUser(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, "session-123", sessionID)
	})

	t.Run("Login invalid user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		user.Email = "invalid-email"

		sessionID, err := service.LoginUser(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, "", sessionID)
	})

	t.Run("Login user not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		repoErr := core_repository_pool.ErrNoRows

		repo.EXPECT().
			FindByEmail(ctx, user).
			Return(core_domain.User{}, repoErr).
			Times(1)

		sessionID, err := service.LoginUser(ctx, user)

		assert.ErrorIs(t, err, core_errors.ErrorUnauthorized)
		assert.Equal(t, "", sessionID)
	})

	t.Run("Login incorrect password", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()

		hash, err := bcrypt.GenerateFromPassword([]byte("another-password"), bcrypt.MinCost)
		if err != nil {
			t.Fatal(err)
		}
		foundUser := user
		foundUser.ID = 12
		foundUser.Password = string(hash)

		repo.EXPECT().
			FindByEmail(ctx, user).
			Return(foundUser, nil).
			Times(1)

		sessionID, err := service.LoginUser(ctx, user)

		assert.ErrorIs(t, err, core_errors.ErrorUnauthorized)
		assert.Equal(t, "", sessionID)
	})

	t.Run("Login Redis error", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		redis := NewMockRedisRepository(t)
		service := NewUserService(repo, redis)
		user := newUser()
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			t.Fatal(err)
		}
		cachedUser := user
		cachedUser.ID = 12
		cachedUser.Password = string(hash)
		redisErr := errors.New("redis unavailable")

		repo.EXPECT().
			FindByEmail(ctx, user).
			Return(cachedUser, nil).
			Times(1)

		redis.EXPECT().
			CreateCache(ctx, cachedUser).
			Return("", redisErr).
			Times(1)

		sessionID, err := service.LoginUser(ctx, user)

		assert.ErrorIs(t, err, redisErr)
		assert.Equal(t, "", sessionID)
	})
}
