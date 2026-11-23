package statsrepo

import (
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
)

type StaticRepository struct {
	pool   core_repository_pool.Pool
	client grpcclient.BudgetService
}

func NewStaticRepository(
	pool core_repository_pool.Pool,
	client grpcclient.BudgetService,
) *StaticRepository {
	return &StaticRepository{
		pool:   pool,
		client: client,
	}
}
