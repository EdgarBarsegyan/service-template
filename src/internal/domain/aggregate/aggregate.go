package aggregate

type EventType string

type Event interface {
	GetEventType() string
}

type Aggregate struct {
	events map[string]Event
}

func New() *Aggregate {
	return &Aggregate{
		events: make(map[string]Event, 0),
	}
}

func (a *Aggregate) TryPublishEvent(event Event) bool {
	_, ok := a.events[event.GetEventType()]
	if ok {
		return false
	}

	a.events[event.GetEventType()] = event
	return true
}

func (a *Aggregate) TryGetEvent(eventType string) (Event, bool) {
	val, ok := a.events[eventType]
	if ok {
		return val, true
	}

	return nil, false
}

func (a *Aggregate) FlashEvents() []Event {
	eventSlice := make([]Event, 0, len(a.events))
	for _, v := range a.events {
		eventSlice = append(eventSlice, v)
	}
	clear(a.events)
	return eventSlice
}
