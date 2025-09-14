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
type UserCreatedEvent struct {
	eventType UserEventType
	Id        Id
}

func NewUserCreatedEvent(id Id) UserCreatedEvent {
	return UserCreatedEvent{
		eventType: UserCreated,
		Id:        id,
	}
}

func (e UserCreatedEvent) GetEventType() string {
	return string(e.eventType)
}

type UserDeletedEvent struct {
	eventType UserEventType
	Id        Id
}

func NewUserDeletedEvent(id Id) UserDeletedEvent {
	return UserDeletedEvent{
		eventType: UserDeleted,
		Id:        id,
	}
}

func (e UserDeletedEvent) GetEventType() string {
	return string(e.eventType)
}

type UserUpdatedEvent struct {
	eventType     UserEventType
	updatedEvents []UserUpdatedEventType
	Id            Id
}

func NewUserUpdatedEvent(id Id) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		eventType:     UserUpdated,
		updatedEvents: make([]UserUpdatedEventType, 0),
		Id:            id,
	}
}

func (e *UserUpdatedEvent) GetEventType() string {
	return string(e.eventType)
}

func (e *UserUpdatedEvent) AddUpdatedEventType(eventType UserUpdatedEventType) {
	e.updatedEvents = append(e.updatedEvents, eventType)
}
