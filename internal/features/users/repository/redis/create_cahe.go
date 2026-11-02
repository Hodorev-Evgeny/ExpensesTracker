package feature_user_redis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *RepositoryRedis) CreateCache(ctx context.Context, user core_domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.rbd.GetTimeout())
	defer cancel()

	sid, err := generateSessionID()
	if err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}
	key := "sessionID:" + sid
	cmd := r.rbd.Set(ctx, key, user.ID, r.rbd.GetTTL())
	if cmd.Err() != nil {
		return "", fmt.Errorf("set session id: %w", cmd.Err())
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
