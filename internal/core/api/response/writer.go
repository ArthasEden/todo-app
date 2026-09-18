// Package response содержит инструменты для формирования HTTP-ответов:
//   - ResponseWriter — обёртка для перехвата статус-кода и точка записи JSON/HTML/Error/NoContent ответов
package response

import (
	"log/slog"
	"net/http"
)

// UnInitStatusCode — сигнальное значение: WriteHeader ещё не вызывался.
// Если статус не был установлен явно, Go возвращает 200 — мы имитируем это поведение.
var UnInitStatusCode = -1

// ResponseWriter — обёртка над http.ResponseWriter, которая запоминает статус-код.
// Стандартный http.ResponseWriter не даёт прочитать статус после WriteHeader,
// а он нужен middleware Trace для логирования.
//
// Встраивание http.ResponseWriter даёт все методы «бесплатно»;
// переопределяем только WriteHeader, чтобы перехватить код.
//
// Инкапсулирует логику записи HTTP-ответов и хранит логгер
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	log        *slog.Logger
}

// NewResponseWriter создаёт обёртку с незаписанным статусом.
func NewResponseWriter(w http.ResponseWriter, log *slog.Logger) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     UnInitStatusCode,
		log:            log,
	}
}

// Переопределяем метод WriteHeader перехватывая статус-код и передаём его дальше в оригинальный ResponseWriter.
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.WriteHeader(statusCode)

	w.statusCode = statusCode
}

// GetStatusCode возвращает записанный статус-код.
// Если WriteHeader не вызывался — возвращает 200 (поведение по умолчанию в Go).
func (w *ResponseWriter) GetStatusCode() int {
	if w.statusCode == UnInitStatusCode {
		return http.StatusOK
	}

	return w.statusCode
}
