package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	t.Run("Recover from panic and return 500 status", func(t *testing.T) {
		handler := Recover()(server.TestPanicHandler())

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		_ = handler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %v", rec.Code)
		}

		expectedBody := "Internal Server Error\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("Handle without panic returns status OK", func(t *testing.T) {
		handler := Recover()(server.TestOKHandler())

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		_ = handler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %v", rec.Code)
		}
	})
}
