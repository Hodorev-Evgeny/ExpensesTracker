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
