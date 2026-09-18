package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"todoapp/internal/core/sentinels"
)

// ErrorResponse — стандартная структура тела ответа при ошибке.
//   - Error   — полный текст ошибки (цепочка от обработчика до причины)
//   - Message — краткое человекочитаемое сообщение (что пытался сделать обработчик)
type ErrorResponse struct {
	Error   string `json:"error"   example:"full error text"`
	Message string `json:"message" example:"short human-readable message"`
}

// JSONResponse сериализует responseBody в JSON и записывает в ответ с указанным статус-кодом.
// Content-Type автоматически определяется json.NewEncoder.
func (w *ResponseWriter) JSONResponse(response any, statusCode int) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.log.Error("write HTTP response", slog.Any("error", err))
	}
}

// NoContentResponse отправляет HTTP 204 No Content — используется при успешном DELETE.
func (w *ResponseWriter) NoContentResponse() {
	w.WriteHeader(http.StatusNoContent)
}

// ErrorResponse транслирует core ошибку в HTTP-статус через errors.Is().
//
// Маппинг:
//   - ErrInvalidArgument → 400
//   - ErrNotFound        → 404
//   - ErrConflict        → 409
//   - остальное          → 500
//
// Каждый тип ошибки логируется на соответствующем уровне (Warn/Debug/Error)
//
// По сути формируем ErrorResponse для дальней отправки в JSONResponse()
func (w *ResponseWriter) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...any)
	)

	switch {
	case errors.Is(err, sentinels.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = w.log.Info

	case errors.Is(err, sentinels.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = w.log.Debug

	case errors.Is(err, sentinels.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = w.log.Info

	default:
		statusCode = http.StatusInternalServerError
		logFunc = w.log.Error
	}

	logFunc(msg, err.Error())

	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}
	w.JSONResponse(response, statusCode)
}

// PanicResponse формирует HTTP 500 при перехвате паники.
// Вызывается из middleware Panic — см. internal/core/transport/http/middleware/common.go.
//
// По сути формируем ErrorResponse для дальней отправки в JSONResponse()
func (w *ResponseWriter) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	w.log.Error(msg, slog.Any("error", err))

	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}
	w.JSONResponse(response, statusCode)
}
