package debug

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

func SetupLogger(enabled bool) {
	level := slog.LevelInfo
	if enabled {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

func WithDebug(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, ctxKey{}, enabled)
}

func IsDebug(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKey{}).(bool)
	return v
}
