package feature_category_transport

import core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"

type CategoryDTO struct {
	CategoryName string `json:"category_name" validate:"required,min=3,max=20"`
	User_id      int    `json:"user_id" validate:"required"`
}

type CategoryResponse struct {
	ID           int    `json:"id"`
	CategoryName string `json:"category_name"`
	User_id      int    `json:"user_id"`
}

func DomainFromResponse(category core_domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:           category.ID,
		CategoryName: category.Name,
		User_id:      category.User_ID,
	}
}

func CategoryDTOFromDomain(category CategoryRequest) core_domain.Category {
	return core_domain.CreateUnincelizedCategory(
		category.CategoryName,
		category.User_id,
	)
}
