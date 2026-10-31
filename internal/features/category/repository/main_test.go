//go:build integration

package feature_repositor_category

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

	pool, clean, err := core_testutil.UpTestDataBase(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't setup test database: %w", err))
		return 1
	}
	defer clean()

	TestPool = pool
	code := m.Run()

	return code
}
