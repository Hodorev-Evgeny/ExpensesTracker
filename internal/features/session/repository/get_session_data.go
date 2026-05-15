package feature_repository_session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *SessionRepository) GetSessionData(
	ctx context.Context,
	key string,
) (core_domain.CookieData, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	val, err := r.rbd.Get(ctx, key).Result()
	fmt.Println("val", val)
	if err != nil {
		return core_domain.CookieData{}, err
	}

	var cookie core_domain.CookieData
	if err := json.Unmarshal([]byte(val), &cookie); err != nil {
		return core_domain.CookieData{}, fmt.Errorf("error unmarshalling cookie data: %w", err)
	}

	return cookie, nil
}
