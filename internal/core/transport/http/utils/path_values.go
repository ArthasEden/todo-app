package core_http_utils

import (
	"fmt"
	"net/http"
	"uuid"

	core_errors "github.com/ArthasEden/todo-app/internal/core/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil(), fmt.Errorf(
			"no key='%s' in path values: %w",
			key, core_errors.ErrInvalidArgument,
		)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil(), fmt.Errorf(
			"path value='%s' by key='%s' not a valid integer: %v: %w",
			pathValue, key, err, core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}
