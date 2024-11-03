package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type Limit struct {
	Exceeded  bool
	Limit     int
	Remaining int
	Reset     time.Time
}
type Limiter interface {
	AddAndCheckLimit(r *http.Request) (Limit, error)
}

func RateLimit(limiter Limiter) server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			if r.Method == http.MethodOptions {
				return next(w, r)
			}

			limit, err := limiter.AddAndCheckLimit(r)
			if err != nil {
				return result.InternalServerErr(err.Error())
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limit.Remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(limit.Reset.Unix(), 10))

			if limit.Exceeded {
				return result.TooManyRequests("Rate limit exceeded")
			}

			return next(w, r)
		}
	}
}

type clientLimit struct {
	requestCount int32
	windowStart  atomic.Int64
}

func getClientID(r *http.Request) string {
	return r.RemoteAddr
}
