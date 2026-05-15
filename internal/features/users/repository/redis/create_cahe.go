package feature_user_redis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *RepositoryRedis) CreateCache(ctx context.Context, user core_domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	jsonBytes, err := json.Marshal(core_domain.CookieData{
		UserID: user.ID,
	})
	if err != nil {
		return "", fmt.Errorf("marshal cookie: %w", err)
	}

	sid, err := generateSessionID()
	if err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}
	key := "sessionID:" + sid
	cmd := r.rbd.Set(ctx, key, string(jsonBytes), time.Hour)
	if cmd.Err() != nil {
		return "", fmt.Errorf("set session id: %w", cmd)
	}

	return sid, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
