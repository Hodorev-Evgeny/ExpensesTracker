package core_repository_pool

import (
	"context"
	"time"
)

type IsolationLevel string

const (
	ReadCommitted  IsolationLevel = "read committed"
	RepeatableRead IsolationLevel = "repeatable read"
	Serializable   IsolationLevel = "serializable"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	GetTimeout() time.Duration
	BeginTx(ctx context.Context, opts TxOptions) (Tx, error)
	Ping(ctx context.Context) error
	Close()
}

type Tx interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type Row interface {
	Scan(dest ...any) error
}

type TxOptions struct {
	IsolationLevel IsolationLevel
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}
