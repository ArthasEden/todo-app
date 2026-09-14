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
