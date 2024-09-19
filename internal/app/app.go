package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlyKarlik/auth-service/internal/config"
	"github.com/FlyKarlik/auth-service/internal/gateways/http/handler"
	"github.com/FlyKarlik/auth-service/internal/gateways/http/server"
	"github.com/FlyKarlik/auth-service/internal/repository"
	"github.com/FlyKarlik/auth-service/internal/tokens"
	"github.com/FlyKarlik/auth-service/internal/usecase"
	"github.com/FlyKarlik/auth-service/pkg/database"
	"github.com/FlyKarlik/auth-service/pkg/logger"
	"github.com/FlyKarlik/auth-service/pkg/tracer"
	_ "github.com/lib/pq"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

// @title auth-service
// @version 1.0
// @description API auth-service

// @host localhost:3000
// @BasePath /
func (a *App) Run() error {
	log := logger.NewLogger(a.cfg.LogLevel)

	log.Debugf("Application initilize...")

	db, err := database.ConnectionPostgresSQLX(a.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect with postgresql error: %w", err)
	}

	tracer, closer, err := tracer.NewJaegerTracer(a.cfg.ServiceName, a.cfg.JaegerHost)
	if err != nil {
		return fmt.Errorf("connect with jaeger tracer error: %w", err)
	}

	repo := repository.New(db, log)
	tokens := tokens.New(a.cfg, repo, log)
	usecase := usecase.New(tokens)
	handler := handler.New(usecase, log, tracer)
	server := server.NewHTTPServer(a.cfg, handler)

	go func() {
		log.Debugf("Starting http server...")
		if err := server.StartHTTPServer(); err != nil {
			log.Errorf("Error starting http server: %s", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Debugf("Application Shutting Down...")

	if err := server.Shuttdown(context.Background()); err != nil {
		return fmt.Errorf("error shuttdowning application: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("error database close connection: %w", err)
	}

	if err := closer.Close(); err != nil {
		return fmt.Errorf("error jaeger tracer close connection: %w", err)
	}

	return nil
}
