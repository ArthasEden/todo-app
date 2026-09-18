package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var usersModel []UserModel

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	LIMIT $1 OFFSET $2
	`

	// Если в запрос приходят параметры nil (limit, offset), то они не будут учитываться, что удобно
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users:%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userModel UserModel

		if err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		); err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		usersModel = append(usersModel, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return userDomainsFromModels(usersModel), nil
}
