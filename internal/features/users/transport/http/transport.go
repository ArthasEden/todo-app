package user_transport_http

import (
	"context"
	"net/http"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_http_server "github.com/ArthasEden/todo-app/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewUserHTTPHandler(userService UserService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userService: userService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
