package logger

import (
	"context"
	"log/slog"
)

type logCtxKey struct{}

// ToContext кладёт логгер в контекст. Вызывается в middleware Logger,
// чтобы все последующие обработчики могли получить логгер с request_id.
func ToContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(
		ctx, logCtxKey{}, log,
	)
}

// FromContext извлекает логгер из контекста.
// Паникует, если логгер не был добавлен — это программная ошибка,
// означающая, что middleware Logger не был подключён.
func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(logCtxKey{}).(*slog.Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}
