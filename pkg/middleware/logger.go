package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()
			log.Printf("Started %s %s", r.Method, r.URL.Path)

			err := next(w, r)

			duration := time.Since(start)
			log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)

			return err
		}
	}
}
