package request

import (
	"fmt"
	"net/http"
	"strconv"
	"todoapp/internal/core/sentinels"
	"uuid"
)

// GetUUIDPathValue извлекает переменную пути (path variable) по ключу key
// и парсит её как UUID. Пример: для маршрута /tasks/{id} ключ — "id".
func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	// Получаем ключ из пути
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.UUID{}, fmt.Errorf("no key='%s' in path values: %w",
			key, sentinels.ErrInvalidArgument)
	}

	// Парсим строку в UUID
	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("path value='%s' by key='%s' not a valid uuid: %v: %w",
			pathValue, key, err, sentinels.ErrInvalidArgument)
	}

	return val, nil
}

// GetIntPathValue извлекает переменную пути по ключу key и парсит её как int.
func GetIntPathValue(r *http.Request, key string) (int, error) {
	// Получаем ключ из пути
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key='%s' in path values: %w",
			key, sentinels.ErrInvalidArgument)
	}

	// Парсим строку в int
	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("path value='%s' by key='%s' not a valid integer: %v: %w",
			pathValue, key, err, sentinels.ErrInvalidArgument)
	}

	return val, nil
}
