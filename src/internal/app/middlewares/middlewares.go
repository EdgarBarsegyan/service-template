package middlewares

import (
	"log/slog"
	"net/http"
	"service-template/internal/app/middlewares/logging"
	"service-template/internal/app/middlewares/recovery"
	"service-template/internal/app/middlewares/timeout"
)

func Build(logger *slog.Logger, handler http.Handler) http.Handler {
	builder := &Builder{}
	builder.Use(&recovery.Middleware2{Logger: logger})
	builder.Use(&logging.Middleware2{Logger: logger})
	builder.Use(&timeout.Middleware{})
	return builder.Build(handler)
}

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type Builder struct {
	middlevares []Middleware
}

func (b *Builder) Use(middleware Middleware) *Builder {
	b.middlevares = append(b.middlevares, middleware)
	return b
}

func (b *Builder) Build(handler http.Handler) http.Handler {
	for i := len(b.middlevares) - 1; i >= 0; i-- {
		handler = b.middlevares[i].Handle(handler)
	}

	return handler
}
