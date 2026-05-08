package feature_service_static

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *StaticService) GetStatic(
	ctx context.Context,
	filters core_domain.FiltersStatic,
) (core_domain.Static, error) {
	list_category, err := s.staticRepository.GetStaticCategories(ctx, filters)
	if err != nil {
		return core_domain.Static{}, fmt.Errorf("error get list category:%w", err)
	}
	list_transaction, err := s.staticRepository.GetStaticTransactions(ctx, filters)
	if err != nil {
		return core_domain.Static{}, fmt.Errorf("error get list transaction:%w", err)
	}

	static := core_domain.NewStatic(list_category, list_transaction)
	return *static, nil
}
