package core_domain

import (
	"fmt"
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
		return fmt.Errorf("name is required")
	}

	if c.User_ID == UnincelizedID {
		return fmt.Errorf("user_id is unincelized")
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
		fmt.Println("limit_id:", *data.Limit_id.Value)
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
