package user

import (
	"context"
	"fmt"
	"service-template/internal/domain/aggregate"
	"service-template/internal/domain/user"
)

type CreatedHandler struct {
}

func (u CreatedHandler) Handle(ctx context.Context, event aggregate.Event) error {
	createdEvent, ok := event.(user.CreatedEvent)
	if !ok {
		return fmt.Errorf("event is not user.CreatedEvent")
	}

	fmt.Printf("Start notifying of user creation. Id: %s\n", createdEvent.Id)

	return nil
}
