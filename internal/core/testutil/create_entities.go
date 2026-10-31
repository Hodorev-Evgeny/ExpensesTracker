package core_testutil

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func CreateTestUser(ctx context.Context, pool core_repository_pool.Pool, user core_domain.User) (core_domain.User, error) {
	ctx, close := context.WithTimeout(ctx, pool.GetTimeout())
	defer close()

	quqry := `
	INSERT INTO trackerapp.users (full_name, email, phone_number, password, time_add)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, full_name, email, phone_number, password, time_add;
	`

	row := pool.QueryRow(ctx, quqry,
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
		return core_domain.User{}, err
	}

	return userDomain, nil
}

func CreateTestCategory(ctx context.Context, pool core_repository_pool.Pool, category core_domain.Category) (core_domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO trackerapp.categories (title, user_id)
		VALUES ($1, $2) 
		RETURNING id, title, user_id;`

	row := pool.QueryRow(ctx,
		query,
		category.Name,
		category.User_ID,
	)

	err := row.Scan(&category.ID,
		&category.Name,
		&category.User_ID,
	)
	if err != nil {
		return core_domain.Category{}, err
	}

	return category, nil
}

func CreateTestTransaction(ctx context.Context, pool core_repository_pool.Pool, transaction core_domain.Transaction) (core_domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO trackerapp.transactions (sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, sum, type_transaction, date, category_id, user_id, comments, time_create, time_changes`

	row := pool.QueryRow(ctx, query,
		transaction.Sum,
		transaction.Type,
		transaction.Date,
		transaction.CategoryID,
		transaction.UserID,
		transaction.Comments,
		transaction.TimeCreated,
		transaction.TimeChange,
	)

	var transactionDomain core_domain.Transaction
	err := row.Scan(&transactionDomain.ID,
		&transactionDomain.Sum,
		&transactionDomain.Type,
		&transactionDomain.Date,
		&transactionDomain.CategoryID,
		&transactionDomain.UserID,
		&transactionDomain.Comments,
		&transactionDomain.TimeCreated,
		&transactionDomain.TimeChange,
	)

	if err != nil {
		return core_domain.Transaction{}, err
	}

	return transactionDomain, nil
}
