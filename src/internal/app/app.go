package app

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service-template/internal/app/api/httpserver"
	"service-template/internal/app/api/swagger"
	"service-template/internal/app/logger"
	"service-template/internal/app/middlewares"
	"service-template/internal/config"
	"service-template/pkg/api"
	"syscall"
	"time"
)

//go:embed api/swagger/*
var swaggerUIFiles embed.FS

func serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, _ := api.GetSwagger()
	spec, _ := d.MarshalJSON()
	w.Write(spec)
}

func getSwaggerUIFS() fs.FS {
	subFS, err := fs.Sub(swaggerUIFiles, "api/swagger")
	if err != nil {
		panic(err)
	}
	return subFS
}

func Run() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	mux := http.NewServeMux()

	server := httpserver.HttpServer{}
	serverInteface := api.NewStrictHandler(server, []api.StrictMiddlewareFunc{})
	swagger.RegisterHandlers(mux)
	handlers := api.HandlerFromMux(serverInteface, mux)
	handlers = middlewares.Build(log, handlers)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handlers,
		ReadTimeout:       10 * time.Second, // защита от slowloris
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		// По умолчанию HTTP/2 включён, если используется TLS. Для простоты — без TLS в примере.
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			// Можно передать данные соединения в контекст запроса
			return ctx
		},
	}

	// Фоновый запуск
	errCh := make(chan error, 1)
	go func() {
		log.Info("server is listening", slog.String(logger.ServerAddr, srv.Addr))
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	// Перехват сигналов и graceful shutdown
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
