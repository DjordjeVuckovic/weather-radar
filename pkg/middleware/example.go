package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log"
	"net/http"
)

func Example() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) error {
			log.Println("Example middleware")
			return next(writer, request)
		}
	}
}
