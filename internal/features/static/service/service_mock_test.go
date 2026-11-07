package feature_service_static

import (
	"context"
	"testing"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/go-openapi/testify/v2/assert"
)

func TestServiceStatic(t *testing.T) {
	ctx := context.Background()

	t.Run("Get correct static", func(t *testing.T) {
		repo := NewMockStaticeRepository(t)
		service := NewStaticService(repo)
		filters := core_domain.FiltersStatic{}

		repo.EXPECT().
			GetStaticCategories(ctx, filters).
			Return([]core_domain.Category{}, nil).
			Times(1)

		repo.EXPECT().
			GetStaticTransactions(ctx, filters).
			Return([]core_domain.Transaction{}, nil).
			Times(1)

		result, err := service.GetStatic(ctx, filters)

		assert.Nil(t, err)
		assert.Equal(t, result, core_domain.Static{
			ShareCategory: make(map[string]float64),
		})
	})

	t.Run("Error get list category", func(t *testing.T) {
		repo := NewMockStaticeRepository(t)
		service := NewStaticService(repo)
		filters := core_domain.FiltersStatic{}

		repo.EXPECT().
			GetStaticCategories(ctx, filters).
			Return([]core_domain.Category{}, core_repository_pool.ErrNoRows).
			Times(1)

		result, err := service.GetStatic(ctx, filters)

		assert.ErrorIs(t, err, core_repository_pool.ErrNoRows)
		assert.Equal(t, result, core_domain.Static{})
	})

	t.Run("Error get list transaction", func(t *testing.T) {
		repo := NewMockStaticeRepository(t)
		service := NewStaticService(repo)
		filters := core_domain.FiltersStatic{}

		repo.EXPECT().
			GetStaticCategories(ctx, filters).
			Return([]core_domain.Category{}, nil).
			Times(1)

		repo.EXPECT().
			GetStaticTransactions(ctx, filters).
			Return([]core_domain.Transaction{}, core_repository_pool.ErrNoRows).
			Times(1)

		result, err := service.GetStatic(ctx, filters)

		assert.ErrorIs(t, err, core_repository_pool.ErrNoRows)
		assert.Equal(t, result, core_domain.Static{})
	})
}
