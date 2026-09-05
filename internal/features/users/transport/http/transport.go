package users_transport_http

type UsersHTTPHandler struct {
	usersService UserService
}

type UserService interface{}

func NewUserHTTPHandler(
	usersService UserService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}
