package grpcclient

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type BudgetConfig struct {
	ADDR    string        `envconfig:"ADDR" required:"true"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func getBudgetConfig() (BudgetConfig, error) {
	budgetConf := BudgetConfig{}
	if err := envconfig.Process("BUDGET", &budgetConf); err != nil {
		return BudgetConfig{}, fmt.Errorf("failed to process budget env var: %w", err)
	}

	return budgetConf, nil
}

func MustGetBudgetConfig() BudgetConfig {
	budgetConfig, err := getBudgetConfig()
	if err != nil {
		panic(err)
	}
	return budgetConfig
}
