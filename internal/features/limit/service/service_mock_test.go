package feature_service_limit

import (
	"context"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/stretchr/testify/assert"
)

func TestServiceLimit(t *testing.T) {
	ctx := context.Background()
	userID := 12

	t.Run("Create correct limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		limit := core_domain.Limit{
			ID:          -1,
			Duration:    time.Now().Add(time.Hour),
			AmountLimit: 1010,
		}
		output := limit
		output.ID = 12

		repo.EXPECT().
			CreateLimit(ctx, limit).
			Return(output, nil).
			Times(1)

		result, err := service.CreateLimit(ctx, limit)

		assert.Nil(t, err)
		assert.Equal(t, output, result)
	})

	t.Run("Create invalid limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		limit := core_domain.Limit{
			ID:          -1,
			Duration:    time.Time{},
			AmountLimit: 1010,
		}

		result, err := service.CreateLimit(ctx, limit)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
		assert.Equal(t, result, core_domain.Limit{})
	})

	t.Run("Create limit failed", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)
		repoErr := core_repository_pool.ErrViolationKey

		limit := core_domain.Limit{
			ID:          -1,
			Duration:    time.Now().Add(time.Hour),
			AmountLimit: 1010,
		}

		repo.EXPECT().
			CreateLimit(ctx, limit).
			Return(core_domain.Limit{}, repoErr).
			Times(1)

		result, err := service.CreateLimit(ctx, limit)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, result, core_domain.Limit{})
	})

	t.Run("Delete correct limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		repo.EXPECT().
			DeleteLimit(ctx, 12, userID).
			Return(nil).
			Times(1)

		err := service.DeleteLimit(ctx, 12, userID)
		assert.Nil(t, err)
	})

	t.Run("Delete invalid limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		err := service.DeleteLimit(ctx, core_domain.UnincelizedID, userID)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
	})

	t.Run("Delete limit failed", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)
		repoErr := core_errors.ErrorNotFoud

		repo.EXPECT().
			DeleteLimit(ctx, 12, userID).
			Return(repoErr).
			Times(1)

		err := service.DeleteLimit(ctx, 12, userID)
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("Get correct limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		limit := core_domain.Limit{
			ID:          -1,
			Duration:    time.Now().Add(time.Hour),
			AmountLimit: 1010,
		}

		repo.EXPECT().
			GetLimit(ctx, 12, userID).
			Return(limit, nil).
			Times(1)

		result, err := service.GetLimit(ctx, 12, userID)

		assert.Nil(t, err)
		assert.Equal(t, limit, result)
	})

	t.Run("Get invalid limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		result, err := service.GetLimit(ctx, -1, userID)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
		assert.Equal(t, result, core_domain.Limit{})
	})

	t.Run("Get limit failed", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)
		repoErr := core_repository_pool.ErrNoRows

		repo.EXPECT().
			GetLimit(ctx, 12, userID).
			Return(core_domain.Limit{}, repoErr).
			Times(1)

		result, err := service.GetLimit(ctx, 12, userID)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, result, core_domain.Limit{})
	})

	t.Run("Get correct all limits", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		output := []core_domain.Limit{
			{
				ID:          12,
				Duration:    time.Now().Add(time.Hour),
				AmountLimit: 1010,
			},
			{
				ID:          13,
				Duration:    time.Now().Add(time.Hour),
				AmountLimit: 110,
			},
		}

		repo.EXPECT().
			GetLimits(ctx, new(2), new(1), userID).
			Return(output, nil).
			Times(1)

		result, err := service.GetLimits(ctx, new(2), new(1), userID)

		assert.Nil(t, err)
		assert.Equal(t, output, result)
	})

	t.Run("Try get limits with invalid offest", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		result, err := service.GetLimits(ctx, new(2), new(-1), userID)

		assert.ErrorIs(t, err, core_errors.ErrorValidation)
		assert.Equal(t, result, []core_domain.Limit{})
	})

	t.Run("Get all limits failed", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)
		repoErr := core_repository_pool.ErrNoRows

		repo.EXPECT().
			GetLimits(ctx, new(1), new(1), userID).
			Return(nil, repoErr).
			Times(1)

		result, err := service.GetLimits(ctx, new(1), new(1), userID)
		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, result, []core_domain.Limit{})
	})

	t.Run("Correct patch limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		limit := core_domain.Limit{
			ID:          12,
			Duration:    time.Now().Add(time.Hour),
			AmountLimit: 1010,
		}
		patch := core_domain.PatchLimit{
			Amount_limit: struct {
				Value *int
				Set   bool
			}{Value: new(123), Set: true},
		}

		repo.EXPECT().
			GetLimit(ctx, 12, userID).
			Return(limit, nil).
			Times(1)

		output := limit
		output.AmountLimit = 123

		repo.EXPECT().
			PatchLimit(ctx, output, userID).
			Return(output, nil).
			Times(1)

		result, err := service.PatchLimit(ctx, 12, patch, userID)

		assert.Nil(t, err)
		assert.Equal(t, result, output)
	})

	t.Run("Failed patch limit: not found", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)
		repoErr := core_errors.ErrorNotFoud

		repo.EXPECT().
			GetLimit(ctx, 12, userID).
			Return(core_domain.Limit{}, repoErr).
			Times(1)

		result, err := service.PatchLimit(ctx, 12, core_domain.PatchLimit{}, userID)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, result, core_domain.Limit{})
	})

	t.Run("Failed patch limit: not validate", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)

		patch := core_domain.PatchLimit{
			Amount_limit: struct {
				Value *int
				Set   bool
			}{Value: new(-123), Set: true},
		}

		result, err := service.PatchLimit(ctx, 12, patch, userID)

		assert.ErrorIs(t, err, core_errors.ErrorBadRequest)
		assert.Equal(t, result, core_domain.Limit{})
	})

	t.Run("Failed patch limit", func(t *testing.T) {
		repo := NewMockLimitRepository(t)
		service := NewLimitService(repo)
		repoErr := core_errors.ErrorConflict

		limit := core_domain.Limit{
			ID:          12,
			Duration:    time.Now().Add(time.Hour),
			AmountLimit: 1010,
		}
		patch := core_domain.PatchLimit{
			Amount_limit: struct {
				Value *int
				Set   bool
			}{Value: new(123), Set: true},
		}

		repo.EXPECT().
			GetLimit(ctx, 12, userID).
			Return(limit, nil).
			Times(1)

		output := limit
		output.AmountLimit = 123

		repo.EXPECT().
			PatchLimit(ctx, output, userID).
			Return(core_domain.Limit{}, repoErr).
			Times(1)

		result, err := service.PatchLimit(ctx, 12, patch, userID)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, result, core_domain.Limit{})
	})
}
