package core_http_types

import (
	"encoding/json"

	"github.com/ArthasEden/todo-app/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

// UnmarshalJSON реализует encoding/json.Unmarshaler.
// Вызывается json.Decoder когда поле присутствует в теле запроса.
// Факт вызова (Set=true) означает, что HTTP клиент намеренно передал это поле.
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

// ToDomain конвертирует HTTP Nullable в доменный Nullable для передачи в сервис.
func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
