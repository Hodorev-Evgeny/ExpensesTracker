package feature_service_categor

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *CategoryService) GetAllCategories(
	ctx context.Context,
	userId int,
) ([]core_domain.Category, error) {
	categoryList, err := s.repository.GetAllCategories(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("Error get list category: %w", err)
	}

	return categoryList, nil
}
