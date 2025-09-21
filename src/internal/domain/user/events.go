package user

type UserEvent interface {
	Type() string
}

type UserEventType string

const (
	UserCreated UserEventType = "UserCreated"
	UserDeleted UserEventType = "UserDeleted"
	UserUpdated UserEventType = "UserUpdated"
)

type UserUpdatedEventType string

const (
	UserEmailUpdated UserUpdatedEventType = "EmailUpdated"
)

// TODO нужно дату создания и модификации добавить в доменную логику
type CreatedEvent struct {
	eventType UserEventType
	Id        Id
}

func NewCreatedEvent(id Id) CreatedEvent {
	return CreatedEvent{
		eventType: UserCreated,
		Id:        id,
	}
}

func (e CreatedEvent) GetEventType() string {
	return string(e.eventType)
}

type DeletedEvent struct {
	eventType UserEventType
	Id        Id
}

func NewDeletedEvent(id Id) DeletedEvent {
	return DeletedEvent{
		eventType: UserDeleted,
		Id:        id,
	}
}

func (e DeletedEvent) GetEventType() string {
	return string(e.eventType)
}

type UpdatedEvent struct {
	eventType     UserEventType
	updatedEvents []UserUpdatedEventType
	Id            Id
}

func NewUpdatedEvent(id Id) *UpdatedEvent {
	return &UpdatedEvent{
		eventType:     UserUpdated,
		updatedEvents: make([]UserUpdatedEventType, 0),
		Id:            id,
	}
}

func (e *UpdatedEvent) GetEventType() string {
	return string(e.eventType)
}

func (e *UpdatedEvent) AddUpdatedEventType(eventType UserUpdatedEventType) {
	e.updatedEvents = append(e.updatedEvents, eventType)
}
