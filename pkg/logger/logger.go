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
type Attr = slog.Attr

var logg *slog.Logger

func GetLogger() *slog.Logger {
	return logg
}

func Init(config Config) {
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
	logg = slog.New(handler)
}

func Info(msg string, args ...any) {
	logg.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	logg.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	logg.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logg.Error(msg, args...)
}

func Int(key string, value int) Attr {
	return slog.Int(key, value)
}
func String(key string, value string) Attr {
	return slog.String(key, value)
}
func Float64(key string, value float64) Attr {
	return slog.Float64(key, value)
}
func Bool(key string, value bool) Attr {
	return slog.Bool(key, value)
}
