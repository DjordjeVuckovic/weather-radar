package server

import (
	"context"
	"errors"
	"github.com/DjordjeVuckovic/weather-radar/pkg/resp"
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	listenAddr              string
	mux                     *http.ServeMux
	middleware              []MiddlewareFunc
	gracefulShutdownTimeout time.Duration
}

type Option func(*Server)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

const (
	defaultGracefulShutdownTimeout = 10 * time.Second
)

func NewServer(listenAddr string, opts ...Option) *Server {
	server := &Server{
		listenAddr:              listenAddr,
		mux:                     http.NewServeMux(),
		gracefulShutdownTimeout: defaultGracefulShutdownTimeout,
	}
	for _, opt := range opts {
		opt(server)
	}
	server.setupHealthCheck()
	return server
}

func WithGracefulShutdownTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.gracefulShutdownTimeout = d
	}
}

func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := http.ListenAndServe(s.listenAddr, s.mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("shutting down the server")
		}
	}()
	slog.Info("Server start listening...", "port", s.listenAddr)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdownTimeout)
	defer cancel()

	if err := s.shutdown(ctx); err != nil {
		slog.Error("Shutting down the server")
		return err
	}
	slog.Info("Shutdown server...")

	return nil
}

func (s *Server) Use(mw ...MiddlewareFunc) {
	s.middleware = append(s.middleware, mw...)
}

func (s *Server) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) {

	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
}

func (s *Server) GET(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodGet, route, h, mw...)
}

func (s *Server) POST(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodPost, route, h, mw...)
}

func (s *Server) PUT(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodPut, route, h, mw...)
}

func (s *Server) DELETE(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodDelete, route, h, mw...)
}

func (s *Server) PATCH(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodPatch, route, h, mw...)
}

func (s *Server) OPTIONS(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodOptions, route, h, mw...)
}

func (s *Server) HEAD(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodHead, route, h, mw...)
}

func (s *Server) TRACE(route string, h HandlerFunc, mw ...MiddlewareFunc) {
	s.handle(http.MethodTrace, route, h, mw...)
}

func (s *Server) handle(method, route string, h HandlerFunc, mw ...MiddlewareFunc) {
	pattern := method + " " + normalizePathSlash(route)
	s.mux.HandleFunc(pattern, s.wrapMiddleware(h, mw...))
}

func (s *Server) wrapMiddleware(h HandlerFunc, mw ...MiddlewareFunc) http.HandlerFunc {
	mws := append(s.middleware, mw...)

	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return handleError(h)
}

func (s *Server) SetupNotFoundHandler() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Route not found", "path", r.URL.Path)
		http.NotFound(w, r)
	})
}

func (s *Server) setupHealthCheck() {
	s.GET("/healthz", func(w http.ResponseWriter, r *http.Request) error {
		err := resp.WriteJSON(w, http.StatusOK, "OK")
		if err != nil {
			return err
		}
		return nil
	})

	s.GET("/ready", func(w http.ResponseWriter, r *http.Request) error {
		err := resp.WriteJSON(w, http.StatusOK, "OK")
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Server) shutdown(ctx context.Context) error {
	return nil
}

func handleError(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		var problem *result.Err
		ok := errors.As(err, &problem)
		if !ok {
			if err != nil {
				_ = resp.WriteProblemJSON(
					w,
					result.NewErr(
						http.StatusInternalServerError,
						err.Error(),
					),
				)
			}
			return
		}
		if problem != nil {
			switch problem.Status {
			case http.StatusNotFound:
				_ = resp.WriteProblemJSON(w, problem)
			case http.StatusBadRequest:
				_ = resp.WriteProblemJSON(w, problem)
			case http.StatusConflict:
				_ = resp.WriteProblemJSON(w, problem)
			case http.StatusUnauthorized:
				_ = resp.WriteProblemJSON(w, problem)
			default:
				_ = resp.WriteProblemJSON(w, problem)
			}
		}
	}
}
