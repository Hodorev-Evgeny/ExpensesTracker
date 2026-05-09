package feature_service_categor

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type CategoryRepository interface {
	CreateCategory(
		ctx context.Context,
		category core_domain.Category,
	) (core_domain.Category, error)

	GetCategoryByID(
		ctx context.Context,
		id int,
	) (core_domain.Category, error)

	GetAllCategories(
		ctx context.Context,
	) ([]core_domain.Category, error)

	RenameCategory(
		ctx context.Context,
		id int,
		category core_domain.Category,
	) (core_domain.Category, error)

	DeleteCategory(
		ctx context.Context,
		id int,
	) error
}

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(
	repository CategoryRepository,
) *CategoryService {
	return &CategoryService{
		repository: repository,
	}
}
