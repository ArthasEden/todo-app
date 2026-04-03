package user_transport_http

import (
	"net/http"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_logger "github.com/ArthasEden/todo-app/internal/core/logger"
	core_http_request "github.com/ArthasEden/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/ArthasEden/todo-app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"min=10,max=15,startswith=+"`
}

type CreateUserResponce struct {
	ID          int    `json:"id"`
	Vesrion     int    `json:"version"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.userService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := dtoFromDomain(userDomain)
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}

func dtoFromDomain(user domain.User) CreateUserResponce {
	return CreateUserResponce{
		ID:      user.ID,
		Vesrion: user.Vesrion,

		FullName:    user.FullName,
		PhoneNumber: user.FullName,
	}
}
