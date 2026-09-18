package users_postgres_repository

import (
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

type UserModel struct {
	ID          uuid.UUID
	Version     int
	FullName    string
	PhoneNumber *string
}

func userModelToDomain(user UserModel) domain.User {
	return domain.User{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func userDomainsFromModels(users []UserModel) []domain.User {
	domainUsers := make([]domain.User, len(users))

	for i, user := range users {
		domainUsers[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return domainUsers
}
