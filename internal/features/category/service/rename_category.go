package feature_service_categor

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *CategoryService) RenameCategory(
	ctx context.Context,
	id int,
	patch core_domain.CategoryUpdate,
) (core_domain.Category, error) {
	category, err := s.GetCategoryByID(ctx, id)
	if err != nil {
		return core_domain.Category{}, fmt.Errorf("error getting category: %w", err)
	}

	if err := category.Update(patch); err != nil {
		return category, fmt.Errorf("error updating category: %w", err)
	}

	categoryUpdate, err := s.repository.RenameCategory(ctx, id, category)
	if err != nil {
		return core_domain.Category{}, fmt.Errorf("error renaming category: %w", err)
	}

	return categoryUpdate, nil
}
