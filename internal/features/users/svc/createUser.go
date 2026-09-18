package users_service

import (
	"context"
	"fmt"
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

func (s *UserService) CreateUser(
	ctx context.Context,
	fullName string,
	phoneNumber *string,
) (domain.User, error) {
	user := domain.CreateUser(
		fullName,
		phoneNumber,
	)

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("Validate user domain: %w", err)
	}

	user.ID = uuid.New()

	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user from repository: %w", err)
	}

	return user, nil
}
