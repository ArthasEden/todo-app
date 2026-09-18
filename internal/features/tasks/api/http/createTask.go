package tasks_transport

import (
	"net/http"
	"time"
	"uuid"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_logger "github.com/ArthasEden/todo-app/internal/core/logger"
	core_http_request "github.com/ArthasEden/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/ArthasEden/todo-app/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string    `json:"title" validate:"required,min=1,max=100"`
	Description  *string   `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID uuid.UUID `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse struct {
	ID           uuid.UUID  `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID uuid.UUID  `json:"author_user_id"`
}

func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		log    = core_logger.FromContext(ctx)
		resp   = core_http_response.NewHTTPResponseHandler(log, w)
		dtoReq = CreateTaskRequest{}
	)

	if err := core_http_request.DecodeAndValidateRequest(r, &dtoReq); err != nil {
		resp.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		dtoReq.Title,
		dtoReq.Description,
		dtoReq.AuthorUserID,
	)

	task, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		resp.ErrorResponse(err, "failed to create task")

		return
	}

	resp.JSONResponse(CreateTaskResponse(task), http.StatusCreated)
}

//func taskDTOFromDomain(task domain.Task) CreateTaskResponse {}
