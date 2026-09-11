package users_postgres_repository

import (
	"context"
	"fmt"
	"uuid"

	core_errors "github.com/ArthasEden/todo-app/internal/core/errors"
)

func (r *UsersRepository) DelUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
