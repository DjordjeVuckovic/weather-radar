package middleware

import (
	"context"
	"fmt"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
)

const CtxFlusherKey = "flusher"

func HTTPStreaming() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			// Set headers for streaming
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Connection", "keep-alive")

			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "Streaming not supported", http.StatusInternalServerError)
				return fmt.Errorf("streaming not supported")
			}

			ctx := context.WithValue(r.Context(), CtxFlusherKey, flusher)
			r = r.WithContext(ctx)

			return next(w, r)
		}
	}
}
