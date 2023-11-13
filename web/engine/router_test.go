package engine_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestSlog(t *testing.T) {
	jsonLoger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	jsonLoger.With("1q1q1", "w2w2w2").Log(context.Background(), slog.LevelInfo, "hhh||error=%v", "111", "2")
	jsonLoger.Log(context.Background(), slog.LevelInfo, "21212")
}
