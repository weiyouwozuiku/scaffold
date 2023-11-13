package engine

import (
	"context"
	"log/slog"
)

type Logger struct {
	slog.Logger
}

func (l *Logger) Log(ctx context.Context)
