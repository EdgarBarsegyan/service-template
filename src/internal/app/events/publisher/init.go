package publisher

import (
	"context"
	"log/slog"
	userEventHandlers "service-template/internal/app/events/handlers/notifications/user"
	"service-template/internal/domain/aggregate"
	domainUser "service-template/internal/domain/user"
	"sync"
)

var handlersInitLock sync.Once
var workersInitLock sync.Once
var handlers map[string][]EventHandler

func Init(workerCount int, logger *slog.Logger) {
	handlersInitLock.Do(func() {
		handlers = map[string][]EventHandler{
			string(domainUser.UserCreated): {
				userEventHandlers.CreatedHandler{},
			},
		}
	})
	initWorkers(workerCount, logger)
}

func initWorkers(count int, logger *slog.Logger) {
	workersInitLock.Do(func() {
		workers := createWorkers(count, logger)
		workerCn := make(chan *worker, count)
		for _, worker := range workers {
			workerCn <- worker
		}
		pl = &pull{
			workerCn: workerCn,
		}
	})
}

type EventHandler interface {
	Handle(ctx context.Context, event aggregate.Event) error
}
