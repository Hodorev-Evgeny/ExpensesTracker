package feature_service_categor

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *CategoryService) CreateCategory(
	ctx context.Context,
	category core_domain.Category,
) (core_domain.Category, error) {
	if err := category.Validate(); err != nil {
		return core_domain.Category{}, fmt.Errorf("invalid category: %w", err)
	}

	categoryDomain, err := s.repository.CreateCategory(ctx, category)
	if err != nil {
		return core_domain.Category{}, fmt.Errorf("error create category: %w", err)
	}

	return categoryDomain, nil
}
