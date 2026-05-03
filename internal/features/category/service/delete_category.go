package feature_service_categor

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *CategoryService) DeleteCategory(
	ctx context.Context,
	id int,
) error {
	if id == core_domain.UnincelizedID {
		return fmt.Errorf("id category is invalid")
	}

	status := s.repository.DeleteCategory(ctx, id)

	return status
}
