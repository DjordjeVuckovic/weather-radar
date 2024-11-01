package server

import (
	"context"
	"errors"
	"github.com/DjordjeVuckovic/weather-radar/pkg/response"
	results "github.com/DjordjeVuckovic/weather-radar/pkg/result"
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

type OptFunc func(*Server)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

const (
	defaultGracefulShutdownTimeout = 10 * time.Second
)

func NewServer(listenAddr string, opts ...OptFunc) *Server {
	server := &Server{
		listenAddr:              listenAddr,
		mux:                     http.NewServeMux(),
		gracefulShutdownTimeout: defaultGracefulShutdownTimeout,
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

func WithGracefulShutdownTimeout(d time.Duration) OptFunc {
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
		slog.Error("shutting down the server")
		return err
	}
	slog.Info("shutdown server...")

	return nil
}

func (s *Server) Use(mw ...MiddlewareFunc) {
	s.middleware = append(s.middleware, mw...)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {

	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
}

func (s *Server) handle(pattern string, handler HandlerFunc, mw ...MiddlewareFunc) {
	s.mux.HandleFunc(pattern, s.wrapMiddleware(handler, mw...))
}

func (s *Server) wrapMiddleware(handler HandlerFunc, mw ...MiddlewareFunc) http.HandlerFunc {
	mws := append(s.middleware, mw...)
	for i := len(mws) - 1; i >= 0; i-- {
		handler = s.middleware[i](handler)
	}

	return handleError(handler)
}

func (s *Server) UseHealthCheck() {
	s.handle("/healthz", func(w http.ResponseWriter, r *http.Request) error {
		err := response.WriteJSON(w, http.StatusOK, "OK")
		if err != nil {
			return err
		}
		return nil
	})

	s.handle("/ready", func(w http.ResponseWriter, r *http.Request) error {
		err := response.WriteJSON(w, http.StatusOK, "OK")
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Server) shutdown(ctx context.Context) error {
	return nil
}

func handleError(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		var problem *results.Problem
		ok := errors.As(err, &problem)
		if !ok {
			_ = response.WriteProblemJSON(
				w,
				results.NewErr(
					http.StatusInternalServerError,
					"internal server error",
					err.Error(),
				),
			)
			return
		}
		if problem != nil {
			switch problem.Code {
			case http.StatusNotFound:
				_ = response.WriteProblemJSON(w, problem)
			case http.StatusBadRequest:
				_ = response.WriteProblemJSON(w, problem)
			case http.StatusConflict:
				_ = response.WriteProblemJSON(w, problem)
			case http.StatusUnauthorized:
				_ = response.WriteProblemJSON(w, problem)
			default:
				_ = response.WriteProblemJSON(w, problem)
			}
		}
	}
}
