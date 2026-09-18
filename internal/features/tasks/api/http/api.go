package tasks_transport

import (
	"context"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_http_server "github.com/ArthasEden/todo-app/internal/core/transport/http/server"
)

type TaskService interface {
	CreateTask(
		context.Context,
		domain.Task,
	) (domain.Task, error)
}

type TasksHTTPHandler struct {
	tasksService TaskService
}

func NewUserHTTPHandler(
	tasksService TaskService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{},
	}
}
