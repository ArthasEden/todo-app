package users_transport_http

import (
	"net/http"

	core_logger "github.com/ArthasEden/todo-app/internal/core/logger"
	core_http_response "github.com/ArthasEden/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/ArthasEden/todo-app/internal/core/transport/http/utils"
)

type DelUserResponse UserDTOResponse

func (h *UsersHTTPHandler) DelUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	respHandler := core_http_response.NewHTTPResponseHandler(log, w)

	id, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	if err := h.usersService.DelUser(ctx, id); err != nil {
		respHandler.ErrorResponse(
			err,
			"failed to delete user",
		)

		return
	}

	respHandler.NoContentReponse()
}
