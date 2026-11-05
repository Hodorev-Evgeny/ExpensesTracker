package core_testutil

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	core_pgx_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func UpTestDataBase(ctx context.Context) (core_repository_pool.Pool, func(), error) {
	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("postgresql"),
		postgres.WithPassword("postgresql"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start postgres: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to PostgreSQL database: %w", err)
	}

	parseconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse database DSN: %w", err)
	}

	testPool, err := pgxpool.NewWithConfig(ctx, parseconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create test pool: %w", err)
	}

	if err := testPool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("could not ping test pool: %w", err)
	}

	root, err := findProjectRoot()
	if err != nil {
		return nil, nil, fmt.Errorf("could not find project root: %w", err)
	}

	migrateFiles, err := filepath.Glob(
		filepath.Join(root, "migrations", "*.up.sql"),
	)
	if err != nil || len(migrateFiles) == 0 {
		return nil, nil, fmt.Errorf("could not find migration files")
	}

	sort.Strings(migrateFiles)
	for _, path := range migrateFiles {
		query, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("could not read migration file: %w", err)
		}

		_, err = testPool.Exec(
			ctx,
			string(query),
			pgx.QueryExecModeSimpleProtocol,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("could apply migration: %w", err)
		}
	}

	cleanup := func() {
		testPool.Close()

		if err := testcontainers.TerminateContainer(container); err != nil {
			log.Printf("Failed to stop PostgreSQL: %v", err)
		} else {
			log.Println("PostgreSQL container stopped")
		}
	}

	pool := core_pgx_pool.NewPool(testPool, 5*time.Second)
	return pool, cleanup, nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
