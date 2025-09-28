package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime/debug"
	"service-template/internal/app/logger"
	"service-template/internal/domain/aggregate"
)

var pl *pull

type worker struct {
	id     int
	logger *slog.Logger
}

func (w *worker) run(ctx context.Context, handler EventHandler, event aggregate.Event) {
	defer func() {
		err := recover()
		if err != nil {
			message, errorType := GetErrorInfo(err)
			jsonData, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				panic(err)
			}

			w.logger.Error(
				"panic",
				slog.Int(EventWorkerId, w.id),
				slog.String(EventType, event.GetEventType()),
				slog.String(EventPayload, string(jsonData)),
				slog.String(logger.StackTrace, string(debug.Stack())),
				slog.String(logger.StackTrace, string(debug.Stack())),
				slog.String(logger.ErrorMessage, message),
				slog.String(logger.ErrorType, errorType),
			)
		}
	}()

	handler.Handle(ctx, event)
}

func createWorkers(count int, logger *slog.Logger) []*worker {
	workers := make([]*worker, 0, count)
	for i := 1; i <= count; i++ {
		workers = append(
			workers,
			&worker{id: i, logger: logger},
		)
	}
	return workers
}

type pull struct {
	workerCn chan *worker
}

func GetPull() *pull {
	if pl == nil {
		panic("event worker pull is not init")
	}

	return pl
}

func (p *pull) containsWorker() bool {
	return len(p.workerCn) > 0
}

func (p *pull) process(ctx context.Context, handler EventHandler, event aggregate.Event) {
	w := <-p.workerCn

	go func() {
		defer func() {
			pl.workerCn <- w
		}()
		defer func() {
			err := recover()
			if err != nil {
				message, errorType := GetErrorInfo(err)
				w.logger.Error(
					"panic",
					slog.Int(EventWorkerId, w.id),
					slog.String(EventType, event.GetEventType()),
					slog.String(logger.StackTrace, string(debug.Stack())),
					slog.String(logger.ErrorMessage, message),
					slog.String(logger.ErrorType, errorType),
				)
			}
		}()

		w.run(ctx, handler, event)
	}()
}

func GetErrorInfo(err any) (message string, errorType string) {
	switch v := err.(type) {
	case error:
		message = v.Error()
		errorType = fmt.Sprintf("%T", v)
	case string:
		message = v
		errorType = "string"
	default:
		message = fmt.Sprintf("%v", v)
		errorType = fmt.Sprintf("%T", v)
	}
	return
}
