package users_service

import (
	"context"
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

func (s UserService) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	userPatch domain.UserPatch,
) (domain.User, error) {
	return domain.User{}, nil
}
