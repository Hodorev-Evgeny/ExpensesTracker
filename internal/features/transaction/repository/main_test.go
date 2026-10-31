//go:build integration

package feature_repository_transaction

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	core_testutil "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/testutil"
)

var TestPool core_repository_pool.Pool

func TestMain(m *testing.M) {
	os.Exit(runTestMain(m))
}

func runTestMain(m *testing.M) int {
	ctx := context.Background()

	pool, clear, err := core_testutil.UpTestDataBase(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("could not connect to database: %w", err))
		return 1
	}
	defer clear()

	TestPool = pool
	code := m.Run()

	return code
}
