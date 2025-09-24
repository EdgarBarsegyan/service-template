package publisher

import (
	"context"
	"fmt"
	userEventHandlers "service-template/internal/app/events/handlers/notifications/user"
	"service-template/internal/domain/aggregate"
	domainUser "service-template/internal/domain/user"
)

var handlers map[string][]EventHandler

func Init() {
	if handlers != nil {
		panic(fmt.Errorf("handlers have already init"))
	}

	handlers = map[string][]EventHandler{
		string(domainUser.UserCreated): {
			userEventHandlers.CreatedHandler{},
		},
	}
}

type EventHandler interface {
	Handle(ctx context.Context, event aggregate.Event) error
}
