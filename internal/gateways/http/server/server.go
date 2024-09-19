package server

import (
	"context"
	"net/http"
	"time"

	"github.com/FlyKarlik/auth-service/internal/config"
	"github.com/FlyKarlik/auth-service/internal/gateways/http/handler"
)

type HTTPServer struct {
	cfg        *config.Config
	handler    *handler.Handler
	httpserver *http.Server
}

func NewHTTPServer(cfg *config.Config, handler *handler.Handler) *HTTPServer {
	return &HTTPServer{cfg: cfg, handler: handler}
}

func (h *HTTPServer) StartHTTPServer() error {

	router := initRoutes(h.handler)

	h.httpserver = &http.Server{
		Addr:              h.cfg.ServerHost,
		MaxHeaderBytes:    1 << 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		Handler:           router,
	}

	return h.httpserver.ListenAndServe()
}

func (s *HTTPServer) Shuttdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return s.httpserver.Shutdown(ctx)
}
