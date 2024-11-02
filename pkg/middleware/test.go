package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log"
	"net/http"
)

func Test() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) error {
			log.Println("Test middleware")
			return next(writer, request)
		}
	}
}
