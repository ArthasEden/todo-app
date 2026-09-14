package logger

import (
	"context"
	"log/slog"
)

type logCtxKey struct{}

func ToContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(
		ctx, logCtxKey{}, log,
	)
}

func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(logCtxKey{}).(*slog.Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}
