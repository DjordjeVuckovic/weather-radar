package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFixedWindowLimiter(t *testing.T) {
	testAddr := "127.0.0.1:1312"

	t.Run("Allow requests within limit", func(t *testing.T) {
		handler := newTestFwLimiterHandler(3*time.Second, 2)

		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = testAddr
			rec := httptest.NewRecorder()

			_ = handler(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status OK, got %v", rec.Code)
			}
		}
	})

	t.Run("Deny requests exceeding limit", func(t *testing.T) {
		handler := newTestFwLimiterHandler(3*time.Second, 2)

		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = testAddr
			rec := httptest.NewRecorder()
			_ = handler(rec, req)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = testAddr
		rec := httptest.NewRecorder()

		_ = handler(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status Too Many Requests, got %v", rec.Code)
		}
	})

	t.Run("Reset limit after window expires", func(t *testing.T) {
		handler := newTestFwLimiterHandler(3*time.Second, 2)

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = testAddr
		rec := httptest.NewRecorder()
		_ = handler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rec.Code)
		}

		time.Sleep(3 * time.Second)

		req = httptest.NewRequest("GET", "/", nil)
		rec = httptest.NewRecorder()
		_ = handler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK after window reset, got %v", rec.Code)
		}
	})
}

func newTestFwLimiterHandler(time time.Duration, maxReq int32) server.HandlerFunc {
	limiter := NewFixedWindowLimiter(FixedWindowLimiterConfig{
		Window:      time,
		MaxRequests: maxReq,
	})
	return RateLimit(limiter)(server.TestOKHandler())
}
