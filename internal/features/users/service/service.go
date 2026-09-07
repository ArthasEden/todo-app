package users_service

import (
	"context"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
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
