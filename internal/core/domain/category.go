package core_domain

import (
	"fmt"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

type Category struct {
	ID       int
	Name     string
	User_ID  int
	Limit_id *int
}

type CategoryUpdate struct {
	Name     Nullable[string]
	Limit_id Nullable[int]
}

func (u *CategoryUpdate) Validate() error {
	if u.Name.Set && u.Name.Value == nil {
		return fmt.Errorf("name cannot be nil")
	}

	return nil
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required: %w", core_errors.ErrorBadRequest)
	}

	if c.User_ID == UnincelizedID || c.User_ID <= 0 {
		return fmt.Errorf("user_id is unincelized: %w", core_errors.ErrorBadRequest)
	}

	if c.ID == UnincelizedID || c.ID <= 0 {
		return fmt.Errorf("id is unincelized: %w", core_errors.ErrorBadRequest)
	}

	if c.Limit_id != nil {
		if *c.Limit_id <= 0 {
			return fmt.Errorf("limit_id is required: %w", core_errors.ErrorBadRequest)
		}
	}

	return nil
}

func (c *Category) Update(data CategoryUpdate) error {
	if err := data.Validate(); err != nil {
		return err
	}

	tmp := *c

	if data.Name.Set {
		tmp.Name = *data.Name.Value
	}

	if data.Limit_id.Set {
		tmp.Limit_id = data.Limit_id.Value
	}

	if tmp.Validate() != nil {
		return fmt.Errorf("new category in invalid")
	}

	*c = tmp
	return nil
}

func NewCategory(
	id int,
	name string,
	user_id int,
	limit_id *int,
) Category {
	return Category{
		ID:       id,
		Name:     name,
		User_ID:  user_id,
		Limit_id: limit_id,
	}
}

func CreateUnincelizedCategory(
	name string,
	user_id int,
	limit_id *int,
) Category {
	return NewCategory(UnincelizedID, name, user_id, limit_id)
}

func RequestUpdateFromDomain(title Nullable[string], limit Nullable[int]) CategoryUpdate {
	return CategoryUpdate{
		Name:     title,
		Limit_id: limit,
	}
}
