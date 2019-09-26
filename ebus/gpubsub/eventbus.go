// Copyright (c) 2019 - goevent authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gpubsub

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gapsquare/goevent"

	"cloud.google.com/go/pubsub"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"

	"google.golang.org/api/option"
)

// DefaultQueueSize is the default queue size per handler for publishing events.
var DefaultQueueSize = 10

// EventBus is a google cloud pubsub event bus client that delegates handling of published events
// to all matched registered handlers, in order of registration
type EventBus struct {
	appID        string
	client       *pubsub.Client
	topic        *pubsub.Topic
	registered   map[goevent.EventHandlerType]struct{}
	registeredMu sync.RWMutex
	errCh        chan goevent.EventBusError
	errorQueue   *errorTopicQueue
}

// NewEventBus creates an EventBus with optional GCP settings
func NewEventBus(projectID, appID string, opts ...option.ClientOption) (*EventBus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}

	//Create or get a topic
	name := appID + "_events"
	topic := client.Topic(name)

	if ok, err := topic.Exists(ctx); err != nil {

		return nil, err
	} else if !ok {
		if topic, err = client.CreateTopic(ctx, name); err != nil {
			return nil, err
		}
	}

	errQueue, err := newErrorTopicQueue(projectID, appID, opts...)
	if err != nil {
		return nil, err
	}

	return &EventBus{
		appID:      appID,
		client:     client,
		topic:      topic,
		registered: map[goevent.EventHandlerType]struct{}{},
		errCh:      make(chan goevent.EventBusError, 100),
		errorQueue: errQueue,
	}, nil
}

// Publish implements the Publish method of the goevent.EventBus interface
func (b *EventBus) Publish(ctx context.Context, event goevent.Event) error {
	e := evt{
		Topic:     event.Topic(),
		Version:   event.Version(),
		Timestamp: event.Timestamp(),
	}

	if event.Data() != nil {
		rawData, err := bson.Marshal(event.Data())
		if err != nil {
			return errors.New("could not marshal event data: " + err.Error())
		}

		e.RawData = bson.Raw{Kind: 3, Data: rawData}
	}
	data, err := bson.Marshal(e)
	if err != nil {
		return errors.New("could not marshal event: " + err.Error())
	}

	//publishCtx := context.Background()
	res := b.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	if _, err := res.Get(ctx); err != nil {
		return errors.New("could not publish event: " + err.Error())
	}

	return nil
}

// AddHandler implements AddHandler of the goevent.EventBus interface
func (b *EventBus) AddHandler(m goevent.EventMatcher, h goevent.EventHandler) error {
	sub := b.subscription(m, h, false)
	go b.handle(m, h, sub)
	return nil
}

// Errors implements Errors method of the the goevent.EventBus interface
func (b *EventBus) Errors() <-chan goevent.EventBusError {
	return b.errCh
}

// Check EventMatcher and EventHandler and gets a subscription
func (b *EventBus) subscription(m goevent.EventMatcher, h goevent.EventHandler, subscription bool) *pubsub.Subscription {
	b.registeredMu.Lock()
	defer b.registeredMu.Unlock()

	if m == nil {
		panic("matcher can't be nil")
	}

	if h == nil {
		panic("handler can't be nil")
	}

	if _, ok := b.registered[h.HandlerType()]; ok {
		panic(fmt.Sprintf("multiple registration for %s", h.HandlerType()))
	}

	b.registered[h.HandlerType()] = struct{}{}
	id := string(h.HandlerType())
	if subscription {
		id = fmt.Sprintf("%s=%s", id, uuid.New())
	}

	ctx := context.Background()
	subscriptionID := b.appID + "_" + id
	sub := b.client.Subscription(subscriptionID)
	if ok, err := sub.Exists(ctx); err != nil {
		panic("could not check subscription: " + err.Error())
	} else if !ok {
		if sub, err = b.client.CreateSubscription(ctx, subscriptionID,
			pubsub.SubscriptionConfig{
				Topic:       b.topic,
				AckDeadline: 60 * time.Second,
				Labels:      map[string]string{"env": b.appID},
			},
		); err != nil {
			panic("could not create subscription: " + err.Error())
		}

	}
	return sub
}

// Handles all events comming in on the channel
func (b *EventBus) handle(m goevent.EventMatcher, h goevent.EventHandler, sub *pubsub.Subscription) {
	for {
		ctx := context.Background()
		if err := sub.Receive(ctx, b.handler(m, h)); err != nil {
			select {
			case b.errCh <- goevent.EventBusError{Err: errors.New("could not receive: " + err.Error())}:
			default:
			}

		}
	}
}

func (b *EventBus) handler(m goevent.EventMatcher, h goevent.EventHandler) func(ctx context.Context, msg *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		data := bson.Raw{
			Kind: 3,
			Data: msg.Data,
		}

		var e evt
		if err := data.Unmarshal(&e); err != nil {
			select {
			case b.errCh <- goevent.EventBusError{Err: errors.New("could not unmarshal event: " + err.Error())}:
			default:
			}

			if e := b.errorQueue.publish("none", msg, errors.New("could not unmarshal event: "+err.Error())); e != nil {
				select {
				case b.errCh <- goevent.EventBusError{Err: errors.New("could not move message to error queue: " + e.Error())}:
				default:
				}
			}

			msg.Ack()
			return
		}

		//Create event of the type
		if data, err := goevent.CreateEventData(e.Topic); err == nil {
			if err := e.RawData.Unmarshal(data); err != nil {
				select {
				case b.errCh <- goevent.EventBusError{Err: errors.New("could not unmarshal event data: " + err.Error())}:
				default:
				}
				if e := b.errorQueue.publish(e.Topic, msg, errors.New("could not unmarshal event: "+err.Error())); e != nil {
					select {
					case b.errCh <- goevent.EventBusError{Err: errors.New("could not move message to error queue: " + e.Error())}:
					default:
					}
				}

				msg.Ack()
				return
			}
			e.data = data
			e.RawData = bson.Raw{}
		}

		event := event{evt: e}

		if !m(event) {
			msg.Ack()
			return
		}

		if err := h.HandleEvent(ctx, event); err != nil {
			select {
			case b.errCh <- goevent.EventBusError{Err: fmt.Errorf("could not handle event (%s): %s", h.HandlerType(), err.Error()), Event: event}:
			default:
			}
			if e := b.errorQueue.publish(e.Topic, msg, errors.New("could not unmarshal event: "+err.Error())); e != nil {
				select {
				case b.errCh <- goevent.EventBusError{Err: errors.New("could not move message to error queue: " + e.Error())}:
				default:
				}
			}

			msg.Ack()
			return
		}
		msg.Ack()
	}
}

type evt struct {
	Topic     goevent.EventTopic  `bson:"event_topic"`
	RawData   bson.Raw            `bson:"data,omitempty"`
	data      goevent.EventData   `bson:"-"`
	Timestamp time.Time           `bson:"timestamp"`
	Version   goevent.VersionType `bson:"version"`
}

// event is the private implementation of the goevent.Event interface
type event struct {
	evt
}

// Topic implements Topic of the goevent.Event interface
func (e event) Topic() goevent.EventTopic {
	return e.evt.Topic
}

// Data implements Data method of the goevent.Event interface
func (e event) Data() goevent.EventData {
	return e.evt.data
}

// Timestamp implements Timestamp method of the goevent.Event interface
func (e event) Timestamp() time.Time {
	return e.evt.Timestamp
}

// Version implements Version method of the goevent.Event interface
func (e event) Version() goevent.VersionType {
	return e.evt.Version
}

//String implements String method of the goevent.Event interface
func (e event) String() string {
	return fmt.Sprintf("%s@%d: [%v]", e.evt.Topic, e.evt.Version, e.data)
}
