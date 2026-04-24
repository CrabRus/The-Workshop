package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Env string // "prod" or "dev"
}

func Setup(cfg Config) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	if cfg.Env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
