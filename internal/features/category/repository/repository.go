package feature_repositor_category

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type CategoryRepository struct {
	pool core_repository_pool.Pool
}

func NewCategoryRepository(
	pool core_repository_pool.Pool,
) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}
