package response

import (
	"log/slog"
	"net/http"
)

var UnStatusCode = -1

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	log        *slog.Logger
}

func NewResponseWriter(w http.ResponseWriter, log *slog.Logger) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     UnStatusCode,
		log:            log,
	}
}

// Переопределяем метод WriteHeader записывая код статуса в нашу структуру
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.WriteHeader(statusCode)

	w.statusCode = statusCode
}

func (w *ResponseWriter) GetStatusCode() int {
	if w.statusCode == UnStatusCode {
		return http.StatusOK
	}

	return w.statusCode
}
