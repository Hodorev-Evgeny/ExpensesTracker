package core_domain

import (
	"fmt"
)

type Category struct {
	ID      int
	Name    string
	User_ID int
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

func NewCategory(
	id int,
	name string,
	user_id int,
) Category {
	return Category{
		ID:      id,
		Name:    name,
		User_ID: user_id,
	}
}

func CreateUnincelizedCategory(
	name string,
	user_id int,
) Category {
	return NewCategory(UnincelizedID, name, user_id)
}
