package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Setup struct {
	Logger *slog.Logger
	Port   uint
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

const (
	portDefault uint = 2021
)

func NewServer(setup *Setup) *Server {
	port := portDefault
	if setup.Port > 0 {
		port = setup.Port
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpserver := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	logger := setup.Logger
	if logger == nil {
		slog.Info("logger was not set, using a pre-built logger")
		logger = slog.Default()
	}

	newServer := Server{
		httpServer: httpserver,
		logger:     logger,
	}

	return &newServer
}

// Start starts http web server.
func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		s.logger.Info("shutting down the web server", slog.String("reason", ctx.Err().Error()))
		s.Close()
	}()

	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.logger.Error("unable to start exporter web server", slog.String("error", err.Error()))
	}
}

func (s *Server) Close() {
	err := s.httpServer.Close()
	if err != nil {
		s.logger.Error("closing http server", slog.String("error", err.Error()))
	}
}
