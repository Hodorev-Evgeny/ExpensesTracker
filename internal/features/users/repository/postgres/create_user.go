package features_users_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *UserRepository) AddUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, close := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer close()

	quqry := `
	INSERT INTO trackerapp.users (full_name, email, phone_number, password, time_add)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, full_name, email, phone_number, password, time_add;
	`

	row := r.pool.QueryRow(ctx, quqry,
		user.Full_name,
		user.Email,
		user.Phone_number,
		user.Password,
		user.Time_add)

	var userDomain core_domain.User
	err := row.Scan(
		&userDomain.ID,
		&userDomain.Full_name,
		&userDomain.Email,
		&userDomain.Phone_number,
		&userDomain.Password,
		&userDomain.Time_add)

	if err != nil {
		if errors.Is(err, core_repository_pool.ErrErrConflict) {
			return core_domain.User{}, fmt.Errorf("this user already exists: %w", core_errors.ErrorExists)
		}
		return core_domain.User{}, fmt.Errorf("error inserting user: %w", err)
	}

	return userDomain, nil
}
