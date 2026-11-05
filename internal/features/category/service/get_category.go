package feature_service_categor

import (
	"context"
	"errors"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *CategoryService) GetCategoryByID(
	ctx context.Context,
	id int,
	userId int,
) (core_domain.Category, error) {
	if id == core_domain.UnincelizedID {
		return core_domain.Category{}, errors.New("unincelized category")
	}

	category, err := s.repository.GetCategoryByID(ctx, id, userId)
	if err != nil {
		return core_domain.Category{}, err
	}

	return category, nil
}
