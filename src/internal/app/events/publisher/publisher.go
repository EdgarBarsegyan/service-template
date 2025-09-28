package publisher

import (
	"context"
	"log/slog"
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
	go func() {
		pull := GetPull()
		for _, v := range events {
			eventHandlers, ok := handlers[v.GetEventType()]
			if !ok {
				continue
			}

			for _, handler := range eventHandlers {
				todoCtx := context.TODO()
				pull.process(todoCtx, handler, v)
			}
		}
	}()
}
