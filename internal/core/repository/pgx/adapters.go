package pgx_pool

import (
	"errors"
	pool "todoapp/internal/core/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pool.ErrNoRows
		}

		return err
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}
