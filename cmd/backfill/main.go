package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hodorev-Evgeny/ExpensesTracker"
	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	corelogger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	pgpool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/infrastructure/grpcclient"
	"go.uber.org/zap"
)

func run(ctx context.Context, logger *corelogger.Logger) error {
	logger.Debug("starting initialization pool connection")
	pgconfig := ExpensesTracker.MustPostgresConfig()
	pool := pgpool.CreatePoolMust(ctx, pgconfig)
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("error pinging pool", zap.Error(err))
		return err
	}

	budgetConf := grpcclient.MustGetBudgetConfig()
	budgetClient, err := grpcclient.NewClientBudget(budgetConf)
	if err != nil {
		logger.Error("error initializing budget client", zap.Error(err))
		return err
	}

	tx, err := pool.BeginTx(ctx, core_repository_pool.TxOptions{
		IsolationLevel: core_repository_pool.ReadCommitted,
	})
	if err != nil {
		logger.Error("error starting transaction", zap.Error(err))
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	querySendOutbox := `
	INSERT INTO trackerapp.outbox (transaction_id)
	SELECT id
	FROM trackerapp.transactions
	WHERE type_transaction = 'Expenditure'
	ON CONFLICT (transaction_id) DO NOTHING;
	`

	if _, err := tx.Exec(ctx, querySendOutbox); err != nil {
		logger.Error("error executing transaction", zap.Error(err))
		return err
	}

	queryGetTransactions := `
		SELECT
			t.id,
			t.sum,
			t.type_transaction,
			t.date,
			t.category_id,
			t.user_id,
			t.comments
		FROM trackerapp.outbox o
		JOIN trackerapp.transactions t
			ON t.id = o.transaction_id
		WHERE o.status = 'unsent'
		ORDER BY o.id;
	`

	rows, err := tx.Query(ctx, queryGetTransactions)
	if err != nil {
		logger.Error("error querying transactions", zap.Error(err))
		return err
	}
	defer rows.Close()

	list := make([]core_domain.Transaction, 0)
	for rows.Next() {
		var transaction core_domain.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Sum,
			&transaction.Type,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.UserID,
			&transaction.Comments,
		)
		if err != nil {
			logger.Error("error scanning rows", zap.Error(err))
			return err
		}
		list = append(list, transaction)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		logger.Error("error iterating rows", zap.Error(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("error committing transaction", zap.Error(err))
		return err
	}

	query := `
	UPDATE trackerapp.outbox
	SET status = 'sent'
	WHERE transaction_id = $1
	  AND status = 'unsent';
	`
	for _, transaction := range list {
		ans, err := budgetClient.ApplySpending(ctx, transaction)
		if err != nil {
			logger.Error("error applying spending transaction", zap.Error(err))
			return err
		}

		if _, err := pool.Exec(ctx, query, transaction.ID); err != nil {
			logger.Error("error update status transaction", zap.Error(err))
			return err
		}
		logger.Info("spending transaction", zap.Any("ans", ans))
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	config := corelogger.MustNewConfig()
	logger, err := corelogger.NewLogger(config)
	if err != nil {
		fmt.Printf("Error initializing logger: %v", err)
		return
	}
	defer logger.Close()

	if err := run(ctx, logger); err != nil {
		logger.Error("error send transaction", zap.Error(err))
		logger.Close()
		os.Exit(1)
	}
}
