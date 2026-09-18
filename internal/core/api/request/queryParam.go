package request

import (
	"fmt"
	"net/http"
	"strconv"
	"todoapp/internal/core/sentinels"
	"uuid"
)

// GetUUIDQueryParam читает query-параметр key и парсит его как UUID.
// Возвращает nil (без ошибки) если параметр отсутствует — это означает «фильтр не задан».
func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	// Получаем значение из квери параметров
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	// Парсим строку в UUID
	val, err := uuid.Parse(param)
	if err != nil {
		return nil, fmt.Errorf("param='%s' by key='%s' not a valid uuid: %v: %w",
			param, key, err, sentinels.ErrInvalidArgument)
	}

	return &val, nil
}

// GetIntQueryParam читает query-параметр key и парсит его как int.
// Возвращает nil если параметр отсутствует (пагинация не задана).
func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	// Парсим строку в int
	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("param='%s' by key='%s' not a valid integer: %v: %w",
			param, key, err, sentinels.ErrInvalidArgument)
	}

	return &val, nil
}
