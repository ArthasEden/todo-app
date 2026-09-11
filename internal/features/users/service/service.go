package users_service

import (
	"context"
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	GetUsers(
		ctx context.Context,
		limit, offset *int,
	) ([]domain.User, error)
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
	DelUser(
		ctx context.Context,
		id uuid.UUID,
	) error
	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		user domain.User,
	) (domain.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(
	userRepository UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}
