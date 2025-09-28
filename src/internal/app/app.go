package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service-template/internal/app/api/httpserver"
	"service-template/internal/app/api/swagger"
	"service-template/internal/app/config"
	"service-template/internal/app/events/publisher"
	"service-template/internal/app/logger"
	"service-template/internal/app/middlewares"
	persistenceInfrastructure "service-template/internal/persistence/infrastructure"
	"service-template/pkg/api"
	"syscall"
	"time"
)

func Run() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	mux := http.NewServeMux()

	server := httpserver.New(log)
	serverInterface := api.NewStrictHandler(server, []api.StrictMiddlewareFunc{})
	swagger.RegisterHandlers(mux)
	handlers := api.HandlerFromMux(serverInterface, mux)
	handlers = middlewares.Build(log, handlers)

	persistenceInfrastructure.MustConfigure(cfg)
	publisher.Init(cfg.EventWorkerCount, log)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handlers,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("server is listening", slog.String(logger.ServerAddr, srv.Addr))
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		log.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			log.Error("error during the launch server", slog.String(logger.ErrorMessage, err.Error()))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", slog.String(logger.ErrorMessage, err.Error()))
		_ = srv.Close()
	}
	log.Info("server stopped")
}
