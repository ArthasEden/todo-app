package users_service

import (
	"context"
	"fmt"
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

func (s *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get users from repository: %w", err)
	}

	return user, nil
}
