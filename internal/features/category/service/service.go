package categoryservice

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
)

type CategoryRepository interface {
	CreateCategory(
		ctx context.Context,
		category core_domain.Category,
	) (core_domain.Category, error)

	GetCategoryByID(
		ctx context.Context,
		id int,
		userId int,
	) (core_domain.Category, error)

	GetAllCategories(
		ctx context.Context,
		userId int,
	) ([]core_domain.Category, error)

	RenameCategory(
		ctx context.Context,
		id int,
		userId int,
		category core_domain.Category,
	) (core_domain.Category, error)

	DeleteCategory(
		ctx context.Context,
		id int,
		userID int,
	) ([]core_domain.DeletedTransaction, error)
}

type CategoryService struct {
	repository CategoryRepository
	client     grpcclient.BudgetService
}

func NewCategoryService(
	repository CategoryRepository,
	client grpcclient.BudgetService,
) *CategoryService {
	return &CategoryService{
		repository: repository,
		client:     client,
	}
}
