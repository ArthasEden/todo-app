package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/ArthasEden/todo-app/internal/core/domain"
	core_logger "github.com/ArthasEden/todo-app/internal/core/logger"
	core_http_request "github.com/ArthasEden/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/ArthasEden/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/ArthasEden/todo-app/internal/core/transport/http/types"
	core_http_utils "github.com/ArthasEden/todo-app/internal/core/transport/http/utils"
)

type PatchUserReq struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

type PatchUserRes UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	respHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var req PatchUserReq
	if err := core_http_request.DecodeAndValidateRequest(r, req); err != nil {
		respHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(req)
	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		respHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	resp := PatchUserRes(userDTOFromDomain(userDomain))

	respHandler.JSONResponse(resp, http.StatusOK)

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fileds:\nFullName: '%s'\nPhoneNumber '%s'",
			req.FullName, req.PhoneNumber,
		),
	)

	w.WriteHeader(http.StatusOK)

}

func userPatchFromRequest(req PatchUserReq) domain.UserPatch {
	return domain.UserPatch{
		FullName:    req.FullName.ToDomain(),
		PhoneNumber: req.PhoneNumber.ToDomain(),
	}
}
