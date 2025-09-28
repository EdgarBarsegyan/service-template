package logging

import (
	"log/slog"
	"net/http"
	"service-template/internal/app/logger"
	"service-template/internal/app/utils"
	"strconv"
	"time"
)

type Middleware struct {
	Logger *slog.Logger
}

func (lm *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lm.Logger.Info(
			"request",
			slog.String(logger.HttpRoute, r.URL.RawPath),
			slog.String(logger.HttpParams, r.URL.RawQuery),
			slog.String(logger.HttpHeaders, utils.GetHeaders(r.Header)),
			slog.String(logger.HttpMethod, r.Method),
			slog.String(logger.HttpRoute, r.RequestURI),
		)

		lrw := &utils.CustomResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		dur := time.Since(start)
		lm.Logger.Info(
			"response",
			slog.String(logger.HttpRoute, r.URL.RawPath),
			slog.String(logger.HttpParams, r.URL.RawQuery),
			slog.String(logger.HttpHeaders, utils.GetHeaders(w.Header())),
			slog.String(logger.HttpMethod, r.Method),
			slog.String(logger.HttpRoute, r.RequestURI),
			slog.String(logger.HttpRoute, r.RequestURI),
			slog.String(logger.HttpStatusCode, strconv.Itoa(lrw.StatusCode)),
			slog.String(logger.HttpMethodDuration, dur.String()),
		)
	})
}
