package ebus

import (
	"context"
	"time"

	"github.com/gapsquare/goevent"
)

const (
	// Topic main mock event
	Topic goevent.EventTopic = "Event"

	// TopicOther - other event
	TopicOther goevent.EventTopic = "OtherEvent"
)

// EventData is a mocked event data, useful in testing :)
type EventData struct {
	Content string
}

func registerEventData() {
	goevent.RegisterEventData(Topic, func() goevent.EventData { return &EventData{} })
}

// EventHandler is a mocked goevent.EventHandler, useful in testing.
type EventHandler struct {
	Type   string
	Events []goevent.Event
	Time   time.Time
	Recv   chan goevent.Event
	// Used to simulate errors when publishing.
	Err error
}

var _ = goevent.EventHandler(&EventHandler{})

// NewEventHandler creates a new EventHandler.
func NewEventHandler(handlerType string) *EventHandler {
	return &EventHandler{
		Type:   handlerType,
		Events: []goevent.Event{},
		Recv:   make(chan goevent.Event, 10),
	}
}

func (m *EventHandler) HandlerType() goevent.EventHandlerType {
	return goevent.EventHandlerType(m.Type)
}

func (m *EventHandler) HandleEvent(ctx context.Context, event goevent.Event) error {
	if m.Err != nil {
		return m.Err
	}

	m.Events = append(m.Events, event)
	m.Time = time.Now()
	m.Recv <- event
	return nil
}

// Reset resets the mock data
func (m *EventHandler) Reset() {
	m.Events = []goevent.Event{}
	m.Time = time.Time{}
}

// Wait is a helper to wait some duration until for an event to be handled
func (m *EventHandler) Wait(d time.Duration) bool {
	select {
	case <-m.Recv:
		return true
	case <-time.After(d):
		return false
	}
}

var _ = goevent.EventBus(&EventBus{})

// EventBus is a mocked goevent.EventBus, useful in testing
type EventBus struct {
	Events []goevent.Event
	Err    error
}

// Publish implements PublishEvent method of the goevent.EventBus interface
func (b *EventBus) Publish(ctx context.Context, event goevent.Event) error {
	if b.Err != nil {
		return b.Err
	}

	b.Events = append(b.Events, event)
	return nil
}

// AddHandler implements AddHandler method of the goevent.EventBus interface
func (b *EventBus) AddHandler(m goevent.EventMatcher, h goevent.EventHandler) error {
	return b.Err
}

// Errors implements Errors method of the goevent.EventBus interface
func (b *EventBus) Errors() <-chan goevent.EventBusError {
	return make(chan goevent.EventBusError)
}
