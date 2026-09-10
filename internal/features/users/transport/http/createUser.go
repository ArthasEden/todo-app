package users_transport_http

import (
	"net/http"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_logger "github.com/ArthasEden/todo-app/internal/core/logger"
	core_http_request "github.com/ArthasEden/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/ArthasEden/todo-app/internal/core/transport/http/response"
)

type CreateUserReq struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserRes UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		log    = core_logger.FromContext(ctx)
		resp   = core_http_response.NewHTTPResponseHandler(log, w)
		dtoReq = CreateUserReq{}
	)

	log.Debug("invoce CreateUser handler")

	if err := core_http_request.DecodeAndValidateRequest(r, &dtoReq); err != nil {
		resp.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	user, err := h.usersService.CreateUser(ctx, domainFromDto(dtoReq))
	if err != nil {
		resp.ErrorResponse(err, "failed to create User")
	}

	resp.JSONResponse(CreateUserRes(userDTOFromDomain(user)), http.StatusCreated)
}

func domainFromDto(dto CreateUserReq) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
