package publisher

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"service-template/internal/app/logger"
	"service-template/internal/domain/aggregate"
)

type IEventPublisher interface {
	PublishEvents([]aggregate.Event)
	Flush()
}

type EventDispatcher interface {
	Dispatch([]aggregate.Event)
}

type EventPublisher struct {
	events     []aggregate.Event
	dispatcher EventDispatcher
}

func New(dispatcher EventDispatcher) *EventPublisher {
	return &EventPublisher{
		dispatcher: dispatcher,
	}
}

func (ep *EventPublisher) PublishEvents(events []aggregate.Event) {
	ep.events = append(ep.events, events...)
}

func (ep *EventPublisher) Flash() {
	ep.dispatcher.Dispatch(ep.events)
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

func (med *MemoryEventDispatcher) Dispatch(events []aggregate.Event) {
	for _, v := range events {
		eventHandlers, ok := handlers[v.GetEventType()]
		if !ok {
			continue
		}

		for _, handler := range eventHandlers {
			go run(v, handler, med.logger)
		}
	}
}

func run(event aggregate.Event, handler EventHandler, l *slog.Logger) {
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

	err := handler.Handle(event)
	if err != nil {
		l.Error(
			"error",
			slog.String(logger.LoggerType, logger.LoggerTypeEventDispatcher),
			slog.String(logger.ErrorMessage, fmt.Sprintf("%v", err)),
			slog.String(logger.ErrorType, fmt.Sprintf("%T", err)),
		)
	}

}
