package feature_category_transport

import core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"

type CategoryDTO struct {
	CategoryName string `json:"category_name" validate:"required,min=3,max=20" example:"Medicine"`
	User_id      int    `json:"user_id" validate:"required" example:"1"`
	Limit_id     *int   `json:"limit_id" example:"1"`
}

type CategoryResponse struct {
	ID           int    `json:"id" example:"4"`
	CategoryName string `json:"category_name" example:"Medicine"`
	User_id      int    `json:"user_id" example:"1"`
	Limit_id     *int   `json:"limit_id" example:"1"`
}

func DomainFromResponse(category core_domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:           category.ID,
		CategoryName: category.Name,
		User_id:      category.User_ID,
		Limit_id:     category.Limit_id,
	}
}

func CategoryDTOFromDomain(category CategoryRequest) core_domain.Category {
	return core_domain.CreateUnincelizedCategory(
		category.CategoryName,
		category.User_id,
		category.Limit_id,
	)
}
