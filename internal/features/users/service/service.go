package users_service

import (
	"context"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

// UsersService — сервис пользователей с бизнес-логикой CRUD-операций.
type UsersService struct {
	usersRepository UsersRepository
}

// UsersRepository — интерфейс репозитория пользователей.
// Определён в пакете сервиса по принципу Dependency Inversion.
type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

// NewUsersService создаёт сервис пользователей с внедрённым репозиторием.
func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
