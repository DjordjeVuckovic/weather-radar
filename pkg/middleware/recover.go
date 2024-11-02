package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log/slog"
	"net/http"
)

func Recover() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("recovered from panic", "error", r)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			return next(w, r)
		}
	}
}
