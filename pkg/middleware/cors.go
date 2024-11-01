package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Origin string
}

func CORS(c Config) server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			origin := os.Getenv("CORS_ORIGINS")
			origins := strings.Split(origin, ",")
			if len(origins) == 0 {
				origins = []string{"*"}
			}
			w.Header().Set("Access-Control-Allow-Origin", c.Origin)
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Api-Key")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return nil
			}
			return next(w, r)
		}
	}
}
