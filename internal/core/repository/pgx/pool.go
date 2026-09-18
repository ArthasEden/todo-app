package pgx_pool

import (
	"context"
	"fmt"
	"time"
	pool "todoapp/internal/core/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("can't make pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("can't ping pgxpool: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)

	return pgxRows{rows}, err
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (pool.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)

	return pgxCommandTag{tag}, err
}

func (p *Pool) Close() {
	p.Pool.Close()
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
