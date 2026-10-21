package features_users_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *UserRepository) FindByEmail(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, full_name, email, phone_number, password, time_add
		FROM trackerapp.users
		WHERE email = $1`

	row := r.pool.QueryRow(ctx, query, user.Email)

	var userLogin core_domain.User
	err := row.Scan(
		&userLogin.ID,
		&userLogin.Full_name,
		&userLogin.Email,
		&userLogin.Phone_number,
		&userLogin.Password,
		&userLogin.Time_add)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf("user not in database: %w", core_repository_pool.ErrNoRows)
		}
		return core_domain.User{}, fmt.Errorf("error while querying user: %w", err)
	}

	return userLogin, nil
}
