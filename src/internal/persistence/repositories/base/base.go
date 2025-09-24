package base

import (
	"context"
	"service-template/internal/domain/aggregate"
)

type EventPublisher interface {
	PublishEvents([]aggregate.Event)
	Flash(ctx context.Context)
}

type BaseRepository struct {
	publisher EventPublisher
}

func New(publisher EventPublisher) *BaseRepository {
	return &BaseRepository{
		publisher: publisher,
	}
}

func (br *BaseRepository) FlashEvents(ctx context.Context, aggregate *aggregate.Aggregate) {
	events := aggregate.FlashEvents()
	br.publisher.PublishEvents(events)
	br.publisher.Flash(ctx)
}
