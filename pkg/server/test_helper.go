package server

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/resp"
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"net/http"
)

func TestOKHandler() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return resp.WriteJSON(w, http.StatusOK, "OK")
	}
}

func TestPanicHandler() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		panic("simulated panic")
	}
}

func TestBadReqErrHandler(_ http.ResponseWriter, _ *http.Request) error {
	return result.NewErr(http.StatusBadRequest, "test error")
}

func TestMiddleware(headerKey string, headerValue string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set(headerKey, headerValue)
			return next(w, r)
		}
	}
}
