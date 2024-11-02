package logger

import (
	"log/slog"
	"os"
)

type Level string

const (
	DebugLevel Level = "debug"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	InfoLevel  Level = "info"
)

type Handler string

const (
	Text Handler = "text"
	Json Handler = "json"
)

type Config struct {
	Level   Level
	Handler Handler
}

func InitSlog(config Config) {
	var level slog.Level
	switch config.Level {
	case DebugLevel:
		level = slog.LevelDebug
	case WarnLevel:
		level = slog.LevelWarn
	case ErrorLevel:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	var handler slog.Handler
	switch config.Handler {
	case Json:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	}
	logg := slog.New(handler)
	slog.SetDefault(logg)
}
