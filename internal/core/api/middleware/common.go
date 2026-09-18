package middleware

import (
	"log/slog"
	"net/http"
	"time"
	"todoapp/internal/core/api/response"
	"todoapp/internal/core/logger"
	"uuid"
)

const reqIDHeader = "X-Request-ID"

// RequestID — middleware, обеспечивающий каждый запрос уникальным идентификатором.
// Если клиент передаёт X-Request-ID — используем его (полезно для распределённой трассировки).
// Иначе генерируем новый X-Request-ID.
// Идентификатор добавляется и в заголовок ответа, чтобы клиент мог его использовать.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get(reqIDHeader)

			if reqID == "" {
				reqID = uuid.New().String()
			}

			r.Header.Set(reqIDHeader, reqID)
			w.Header().Set(reqIDHeader, reqID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger — middleware, кладущий логгер в контекст запроса.
// Обогащает логгер полями request_id и url, чтобы все последующие
// обработчики автоматически логировали эти поля.
//
// Важно: этот middleware должен идти ПОСЛЕ RequestID, чтобы request_id уже был доступен.
func Logger(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get(reqIDHeader)

			// Постоянные поля для логгера
			log := log.With(
				slog.String("request_id", reqID),
				slog.String("url", r.URL.String()),
			)

			ctx := logger.ToContext(r.Context(), log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Trace — middleware для логирования входящих запросов и времени их обработки.
// Использует ResponseWriter-обёртку, чтобы перехватить статус-код ответа,
// который иначе недоступен после вызова WriteHeader.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			before := time.Now().UTC()
			rw := response.NewResponseWriter(w, log)

			log.Debug(
				">>> incoming HTTP request",
				"http_method", r.Method,
				"time", before,
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				"status_code", rw.GetStatusCode(),
				"time", time.Since(before),
			)
		})
	}
}

// Panic — middleware для перехвата паник и возврата HTTP 500.
// Без этого middleware паника в обработчике уронила бы всю горутину,
// а стандартная библиотека Go вернула бы пустой ответ клиенту.
//
// Использует defer + recover — стандартный паттерн обработки паник в Go.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())

			defer func() {
				if p := recover(); p != nil {
					log.Error("some panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
