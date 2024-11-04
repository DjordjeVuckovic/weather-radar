package middleware

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	corsConfig := CORSConfig{
		Origin: "https://weather.api.com",
	}

	handler := CORS(corsConfig)(server.TestOKHandler())

	t.Run("Sets CORS headers on GET request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		_ = handler(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != corsConfig.Origin {
			t.Errorf("Expected Access-Control-Allow-Origin to be %s, got %s", corsConfig.Origin, rec.Header().Get("Access-Control-Allow-Origin"))
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rec.Code)
		}
	})

	t.Run("Handles OPTIONS request with No Content status", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		rec := httptest.NewRecorder()

		_ = handler(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status code 204 No Content, got %d", rec.Code)
		}

		if rec.Header().Get("Access-Control-Allow-Origin") != corsConfig.Origin {
			t.Errorf("Expected Access-Control-Allow-Origin to be %s, got %s", corsConfig.Origin, rec.Header().Get("Access-Control-Allow-Origin"))
		}
	})
}
