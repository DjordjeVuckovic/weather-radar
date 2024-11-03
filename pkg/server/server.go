package server

import (
	"context"
	"errors"
	"github.com/DjordjeVuckovic/weather-radar/docs"
	"github.com/DjordjeVuckovic/weather-radar/pkg/resp"
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	listenAddr string
	mux        *http.ServeMux
	httpServer *http.Server
	middleware []MiddlewareFunc

	gracefulShutdownTimeout time.Duration
	ShutdownSig             chan struct{}
}

type Option func(*Server)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

const (
	defaultGracefulShutdownTimeout = 10 * time.Second
)

func NewServer(listenAddr string, opts ...Option) *Server {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	server := &Server{
		listenAddr:              listenAddr,
		mux:                     mux,
		gracefulShutdownTimeout: defaultGracefulShutdownTimeout,
		httpServer:              s,
		ShutdownSig:             make(chan struct{}),
	}

	for _, opt := range opts {
		opt(server)
	}

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
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("shutting down the server")
		}
	}()
	slog.Info("Server start listening...", "port", s.listenAddr)

	<-ctx.Done()

	close(s.ShutdownSig)

	ctx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
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

func (s *Server) SetupNotFoundHandler() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Warn("Route not found", "path", r.URL.Path)
		http.NotFound(w, r)
	})
}

func (s *Server) SetupSwagger() {
	docs.SwaggerInfo.Title = "Weather Radar"
	docs.SwaggerInfo.Description = "Weather Radar API"
	docs.SwaggerInfo.Version = "1.0"
	s.mux.Handle("/swagger-ui/", httpSwagger.WrapHandler)
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
