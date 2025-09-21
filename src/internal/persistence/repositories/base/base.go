package base

import "service-template/internal/domain/aggregate"

type EventPublisher interface {
	PublishEvents([]aggregate.Event)
	Flash()
}

type BaseRepository struct {
	publisher EventPublisher
}

func New(publisher EventPublisher) *BaseRepository {
	return &BaseRepository{
		publisher: publisher,
	}
}

func (br *BaseRepository) FlashEvents(aggregate *aggregate.Aggregate) {
	events := aggregate.FlashEvents()
	br.publisher.PublishEvents(events)
	br.publisher.Flash()
}
