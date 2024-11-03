package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerUseMiddleware(t *testing.T) {
	s := NewServer(":1312")
	s.Use(TestMiddleware("X-Test-Middleware", "true"))

	s.GET("/test", TestOKHandler())

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	s.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	if rec.Header().Get("X-Test-Middleware") != "true" {
		t.Errorf("Expected X-Test-Middleware header to be true")
	}
}

func TestServerUseMiddlewares(t *testing.T) {
	s := NewServer(":1312")

	s.Use(TestMiddleware("X-Test-Middleware", "true"))
	s.Use(TestMiddleware("Content-Type", "application/json"))
	s.Use(TestMiddleware("X-Custom-Header", "custom-value"))

	s.GET("/test", TestOKHandler())

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	s.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	if rec.Header().Get("X-Test-Middleware") != "true" {
		t.Errorf("Expected X-Test-Middleware header to be true")
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json")
	}
	if rec.Header().Get("X-Custom-Header") != "custom-value" {
		t.Errorf("Expected X-Custom-Header to be custom-value")
	}
}

func TestServerWrapMiddlewareOrder(t *testing.T) {
	s := NewServer(":1312")

	var executionOrder []string
	s.Use(func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			executionOrder = append(executionOrder, "global")
			return next(w, r)
		}
	})
	s.GET("/test", func(w http.ResponseWriter, r *http.Request) error {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
		return nil
	}, func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			executionOrder = append(executionOrder, "route")
			return next(w, r)
		}
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	s.mux.ServeHTTP(rec, req)

	expectedOrder := []string{"global", "route", "handler"}
	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Expected %s at position %d, got %s", expected, i, executionOrder[i])
		}
	}
}

func TestServerNotFoundHandler(t *testing.T) {
	s := NewServer(":1312")
	s.SetupNotFoundHandler()

	req := httptest.NewRequest("GET", "/non-existent", nil)
	rec := httptest.NewRecorder()

	s.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", rec.Code)
	}
}

func TestServerMethodHandling(t *testing.T) {
	s := NewServer(":1312")

	s.GET("/get", TestOKHandler())
	s.POST("/post", TestOKHandler())
	s.PUT("/put", TestOKHandler())
	s.DELETE("/delete", TestOKHandler())
	s.PATCH("/patch", TestOKHandler())

	tests := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/get", http.StatusOK},
		{"POST", "/post", http.StatusOK},
		{"PUT", "/put", http.StatusOK},
		{"DELETE", "/delete", http.StatusOK},
		{"PATCH", "/patch", http.StatusOK},
		{"PATCH", "/non-existent", http.StatusNotFound},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		rec := httptest.NewRecorder()
		s.mux.ServeHTTP(rec, req)

		if rec.Code != test.expected {
			t.Errorf("Expected status code %d, got %d for %s %s", test.expected, rec.Code, test.method, test.path)
		}
	}
}

func TestGracefulShutdown(t *testing.T) {
	s := NewServer(":1312", WithGracefulShutdownTimeout(1*time.Second))

	go func() {
		_ = s.Start()
	}()

	close(s.ShutdownSig)

	time.Sleep(500 * time.Millisecond)

	select {
	case <-s.ShutdownSig:
	default:
		t.Error("Expected server to start shutdown")
	}
}

func TestHandleError(t *testing.T) {
	s := NewServer(":1312")
	s.GET("/error", TestBadReqErrHandler)

	req := httptest.NewRequest("GET", "/error", nil)
	rec := httptest.NewRecorder()

	s.mux.ServeHTTP(rec, req)

	// Check the status code for error response
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", rec.Code)
	}

	// Verify response body
	if rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Errorf("Expected Content-Type application/problem+json, got %s", rec.Header().Get("Content-Type"))
	}
}
