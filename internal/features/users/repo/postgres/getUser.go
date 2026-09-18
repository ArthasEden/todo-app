package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_errors "github.com/ArthasEden/todo-app/internal/core/errors"
	core_postgres_pool "github.com/ArthasEden/todo-app/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var userModel UserModel

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	WHERE id=$1
	`

	// Если в запрос приходят параметры nil (limit, offset), то они не будут учитываться, что удобно
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return userModelToDomain(userModel), nil
}
