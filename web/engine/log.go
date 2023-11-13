package engine

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

var Logger logger

type logger struct {
	log *slog.Logger
}

func NewLogger() error {
	Logger.log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return nil
}

type LogFuncWithCtx func(ctx context.Context, prefix string, format string, args ...any)

func (l *logger) Errorf(ctx context.Context, prefix string, format string, args ...any) {
	l.log.Log(ctx, slog.LevelError, prefix, args...)
}
func (l *logger) Infof(ctx context.Context,prefix )