package features_users_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *UserRepository) FindByEmail(ctx context.Context, user core_domain.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id
		FROM trackerapp.users
		WHERE email = $1 AND password = $2`

	row := r.pool.QueryRow(ctx, query, user.Email, user.Password)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return -1, fmt.Errorf("user not in database: %w", core_repository_pool.ErrNoRows)
		}
		return -1, fmt.Errorf("error while querying user: %w", err)
	}

	return userID, nil
}
