package users_service

type UserService struct {
	userRepository UserRepository
}

type UserRepository interface {
}

func NewUserService(userRepository UserRepository) *UserService {
	return &{
		userRepository: userRepository,
	}
}
