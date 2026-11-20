package pgpool

import (
	"context"
	"errors"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IsolationLevel string

const (
	ReadCommitted  IsolationLevel = "read committed"
	RepeatableRead IsolationLevel = "repeatable read"
	Serializable   IsolationLevel = "serializable"
)

type Row struct {
	pgx.Row
}

type Rows struct {
	pgx.Rows
}

type Tx struct {
	pgx.Tx
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

func (t *Tx) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (core_repository_pool.CommandTag, error) {
	tag, err := t.Tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, MappingError(err)
	}

	return CommandTag{CommandTag: tag}, nil
}

func (t *Tx) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) core_repository_pool.Row {
	return Row{Row: t.Tx.QueryRow(ctx, query, args...)}
}

func (t *Tx) Query(
	ctx context.Context,
	query string,
	args ...any,
) (core_repository_pool.Rows, error) {
	rows, err := t.Tx.Query(ctx, query, args...)
	if err != nil {
		return nil, MappingError(err)
	}

	return Rows{Rows: rows}, nil
}

func (r Row) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		newError := MappingError(err)
		return newError
	}
	return nil
}

type CommandTag struct {
	pgconn.CommandTag
}

func MappingError(err error) error {
	var pgxErr *pgconn.PgError

	if errors.Is(err, pgx.ErrNoRows) {
		return core_repository_pool.ErrNoRows
	}

	if errors.As(err, &pgxErr) {
		switch pgxErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return core_repository_pool.ErrViolationKey
		case pgerrcode.UniqueViolation:
			return core_errors.ErrConflict
		}
	}

	return err
}
