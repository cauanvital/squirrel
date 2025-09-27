//go:build go1.8
// +build go1.8

package squirrel

import (
	"context"
	"database/sql"
)

func (d *selectData) ExecContext(ctx context.Context) (sql.Result, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	ctxRunner, ok := d.RunWith.(ExecerContext)
	if !ok {
		return nil, ErrNoContextSupport
	}
	return ExecContextWith(ctx, ctxRunner, d)
}

func (d *selectData) QueryContext(ctx context.Context) (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	ctxRunner, ok := d.RunWith.(QueryerContext)
	if !ok {
		return nil, ErrNoContextSupport
	}
	return QueryContextWith(ctx, ctxRunner, d)
}

func (d *selectData) QueryRowContext(ctx context.Context) RowScanner {
	if d.RunWith == nil {
		return &Row{err: ErrRunnerNotSet}
	}
	queryRower, ok := d.RunWith.(QueryRowerContext)
	if !ok {
		if _, ok := d.RunWith.(QueryerContext); !ok {
			return &Row{err: ErrRunnerNotQueryRunner}
		}
		return &Row{err: ErrNoContextSupport}
	}
	return QueryRowContextWith(ctx, queryRower, d)
}

// ExecContext builds and ExecContexts the query with the Runner set by RunWith.
func (b selectBuilder) ExecContext(ctx context.Context) (sql.Result, error) {
	return b.data.ExecContext(ctx)
}

// QueryContext builds and QueryContexts the query with the Runner set by RunWith.
func (b selectBuilder) QueryContext(ctx context.Context) (*sql.Rows, error) {
	return b.data.QueryContext(ctx)
}

// QueryRowContext builds and QueryRowContexts the query with the Runner set by RunWith.
func (b selectBuilder) QueryRowContext(ctx context.Context) RowScanner {
	return b.data.QueryRowContext(ctx)
}

// ScanContext is a shortcut for QueryRowContext().Scan.
func (b selectBuilder) ScanContext(ctx context.Context, dest ...interface{}) error {
	return b.QueryRowContext(ctx).Scan(dest...)
}
