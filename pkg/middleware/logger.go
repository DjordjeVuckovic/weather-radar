package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log/slog"
	"net/http"
	"time"
)

func Logger() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			start := time.Now()
			slog.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
				slog.String("uri", r.URL.Path),
				slog.String("method", r.Method),
			)
			err := next(w, r)
			duration := time.Since(start)
			slog.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
				slog.String("uri", r.URL.Path),
				slog.String("method", r.Method),
				slog.Duration("duration", duration),
			)

			return err
		}
	}
}
