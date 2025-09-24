package publisher

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"service-template/internal/app/logger"
	"service-template/internal/domain/aggregate"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, events []aggregate.Event)
}

type EventPublisher struct {
	events     []aggregate.Event
	dispatcher Dispatcher
}

func New(dispatcher Dispatcher) *EventPublisher {
	return &EventPublisher{
		dispatcher: dispatcher,
	}
}

func (ep *EventPublisher) PublishEvents(events []aggregate.Event) {
	ep.events = append(ep.events, events...)
}

func (ep *EventPublisher) Flash(ctx context.Context) {
	ep.dispatcher.Dispatch(ctx, ep.events)
	clear(ep.events)
}

type MemoryEventDispatcher struct {
	logger *slog.Logger
}

func NewMemoryEventDispatcher(logger *slog.Logger) *MemoryEventDispatcher {
	return &MemoryEventDispatcher{
		logger: logger,
	}
}

func (med *MemoryEventDispatcher) Dispatch(ctx context.Context, events []aggregate.Event) {
	for _, v := range events {
		eventHandlers, ok := handlers[v.GetEventType()]
		if !ok {
			continue
		}

		for _, handler := range eventHandlers {
			go run(context.TODO(), v, handler, med.logger)
		}
	}
}

func run(ctx context.Context, event aggregate.Event, handler EventHandler, l *slog.Logger) {
	defer func() {
		if err := recover(); err != nil {
			l.Error(
				"panic",
				slog.String(logger.LoggerType, logger.LoggerTypeEventDispatcher),
				slog.String(logger.StackTrace, string(debug.Stack())),
				slog.String(logger.ErrorMessage, fmt.Sprintf("%v", err)),
				slog.String(logger.ErrorType, fmt.Sprintf("%T", err)),
			)
		}
	}()

	err := handler.Handle(ctx, event)
	if err != nil {
		l.Error(
			"error",
			slog.String(logger.LoggerType, logger.LoggerTypeEventDispatcher),
			slog.String(logger.ErrorMessage, fmt.Sprintf("%v", err)),
			slog.String(logger.ErrorType, fmt.Sprintf("%T", err)),
		)
	}

}
