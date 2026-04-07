package users_service

import (
	"context"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

type UsersService struct {
	usersRepository UsersRepository
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
