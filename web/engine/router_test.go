package engine_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestSlog(t *testing.T) {
	jsonLoger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})).With()
	jsonLoger.Log(context.Background(), slog.LevelInfo, "hhh||error=%v", "111", "2")
}
