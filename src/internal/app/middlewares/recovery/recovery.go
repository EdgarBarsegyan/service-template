package recovery

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"service-template/internal/app/logger"
	"service-template/internal/app/utils"
)

type Middleware2 struct {
	Logger *slog.Logger
}

func (rm *Middleware2) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorInfo := getErrorInfo(err)

				rm.Logger.Error(
					"panic",
					slog.String(logger.HttpRoute, r.RequestURI),
					slog.String(logger.HttpParams, r.URL.RawQuery),
					slog.String(logger.HttpHeaders, utils.GetHeaders(r.Header)),
					slog.String(logger.HttpMethod, r.Method),
					slog.String(logger.StackTrace, string(debug.Stack())),
					slog.String(logger.ErrorMessage, errorInfo.message),
					slog.String(logger.ErrorType, errorInfo.message),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type Middleware struct {
	logger *slog.Logger
	next   http.Handler
}

func (rm *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errorInfo := getErrorInfo(err)

			rm.logger.Error(
				"panic",
				slog.String(logger.HttpRoute, r.RequestURI),
				slog.String(logger.HttpParams, r.URL.RawQuery),
				slog.String(logger.HttpHeaders, utils.GetHeaders(r.Header)),
				slog.String(logger.HttpMethod, r.Method),
				slog.String(logger.StackTrace, string(debug.Stack())),
				slog.String(logger.ErrorMessage, errorInfo.message),
				slog.String(logger.ErrorType, errorInfo.message),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	rm.next.ServeHTTP(w, r)
}

type errorInfo struct {
	message   string
	errorType string
}

func getErrorInfo(err any) errorInfo {
	errorInfo := errorInfo{}

	switch v := err.(type) {
	case error:
		errorInfo.message = v.Error()
		errorInfo.errorType = fmt.Sprintf("%T", v)
	case string:
		errorInfo.message = v
		errorInfo.errorType = "string"
	default:
		errorInfo.message = fmt.Sprintf("%v", v)
		errorInfo.errorType = fmt.Sprintf("%T", v)
	}

	return errorInfo
}
