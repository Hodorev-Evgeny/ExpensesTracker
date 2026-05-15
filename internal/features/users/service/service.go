package feature_user_service

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type UserRepository interface {
	AddUser(
		ctx context.Context,
		user core_domain.User,
	) (core_domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]core_domain.User, error)

	ExtraditionUser(
		ctx context.Context,
		id int,
	) (core_domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch core_domain.User,
	) (core_domain.User, error)

	FindByEmail(
		ctx context.Context,
		user core_domain.User,
	) (int, error)
}

type RedisRepository interface {
	CreateCache(
		ctx context.Context,
		user core_domain.User,
	) (string, error)
}

type UserService struct {
	userRepository      UserRepository
	userRedisRepository RedisRepository
}

func NewUserService(
	userRepository UserRepository,
	repository RedisRepository) *UserService {
	return &UserService{
		userRepository:      userRepository,
		userRedisRepository: repository,
	}
}
