// Package request содержит хелперы для чтения данных из HTTP-запроса:
//   - DecodeAndValidateRequest — десериализация тела запроса + валидация
//   - GetUUIDPathValue / GetIntPathValue — извлечение переменных из пути URL
//   - GetUUIDQueryParam / GetIntQueryParam / GetDateQueryParam — query-параметры
//
// Все функции оборачивают ошибки в ErrInvalidArgument, чтобы транспортный слой
// автоматически вернул HTTP 400 Bad Request.
package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoapp/internal/core/sentinels"

	"github.com/go-playground/validator/v10"
)

// requestValidator — глобальный экземпляр валидатора (go-playground/validator).
// Используется для валидации struct-тегов вида `validate:"required,min=1,max=100"`.
// Создаётся один раз на пакет (потокобезопасен).
var requestValidator = validator.New()

// validatable — интерфейс для DTO с собственным методом Validate().
// Если тип реализует этот интерфейс, используется его метод вместо struct-тегов.
// Это нужно для сложных правил валидации, которые нельзя выразить тегами
// (например, Nullable-поля в PATCH-запросах).
type validatable interface {
	Validate() error
}

// DecodeAndValidateRequest десериализует тело запроса из JSON в dest
// и затем проверяет валидность данных.
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", sentinels.ErrInvalidArgument, err)
	}

	var err error

	// Проверяем значение на тип validatable
	v, ok := dest.(validatable)
	if ok {
		// Если dest реализует validatable — вызывается его Validate()
		err = v.Validate()
	} else {
		// Иначе — валидируются struct-теги через go-playground/validator
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, sentinels.ErrInvalidArgument)
	}

	return nil
}
