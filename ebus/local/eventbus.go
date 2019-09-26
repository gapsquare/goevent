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

package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/gapsquare/goevent"

	"github.com/google/uuid"
)

// DefaultQueueSize is the default queue size per handler for publishing events.
var DefaultQueueSize = 10

// EventBus is a local event bus that delegates handling of published events
// to all matching registered handlers
type EventBus struct {
	group        *Group
	registered   map[goevent.EventHandlerType]struct{}
	registeredMu sync.RWMutex
	errCh        chan goevent.EventBusError
	wg           sync.WaitGroup
}

// Group is publishing group shared by multiple event busses locally
type Group struct {
	bus   map[string]chan evt
	busMu sync.RWMutex
}

// NewEventBus creates a local EventBus
func NewEventBus(g *Group) (*EventBus, error) {
	if g == nil {
		g = NewGroup()

	}

	return &EventBus{
		group:      g,
		registered: map[goevent.EventHandlerType]struct{}{},
		errCh:      make(chan goevent.EventBusError, 100),
	}, nil
}

// Publish implements the Publish method of the goevent.EventBus interface
func (b *EventBus) Publish(ctx context.Context, event goevent.Event) error {
	b.group.publish(ctx, event)
	return nil
}

//AddHandler implements the AddHandler of the goevent.EventBus interface
func (b *EventBus) AddHandler(m goevent.EventMatcher, h goevent.EventHandler) error {
	ch := b.channel(m, h, false)
	go b.handle(m, h, ch)
	return nil
}

// Errors implements the Errors method of the goevent.EventBus interface
func (b *EventBus) Errors() <-chan goevent.EventBusError {
	return b.errCh
}

// handle handles all events coming in on the channel
func (b *EventBus) handle(m goevent.EventMatcher, h goevent.EventHandler, ch <-chan evt) {
	b.wg.Add(1)
	defer b.wg.Done()
	for e := range ch {
		if !m(e.event) {
			continue
		}
		if err := h.HandleEvent(e.ctx, e.event); err != nil {
			select {
			case b.errCh <- goevent.EventBusError{Err: fmt.Errorf("could not handle event (%s): %s", h.HandlerType(), err.Error()), Event: e.event}:
			default:
			}
		}
	}
}

// channel checks the matcher and eventhandler and gets the event channel from the group
func (b *EventBus) channel(m goevent.EventMatcher, h goevent.EventHandler, observer bool) <-chan evt {
	b.registeredMu.Lock()
	defer b.registeredMu.Unlock()

	if m == nil {
		panic("matcher can't be nil")
	}

	if h == nil {
		panic("handler can't be nil")
	}

	if _, ok := b.registered[h.HandlerType()]; ok {
		panic(fmt.Sprintf("multiple registrations for %s", h.HandlerType()))
	}

	b.registered[h.HandlerType()] = struct{}{}

	id := string(h.HandlerType())
	if observer {
		id = fmt.Sprintf("%s-%s", id, uuid.New())
	}

	return b.group.channel(id)
}

// Wait for all channels to close in the event bus group
func (b *EventBus) Wait() {
	b.wg.Wait()
}

// Close all the channels int the events bus group
func (b *EventBus) Close() {
	b.group.Close()
}

// NewGroup creates a Group
func NewGroup() *Group {
	return &Group{
		bus: map[string]chan evt{},
	}
}

type evt struct {
	ctx   context.Context
	event goevent.Event
}

func (g *Group) channel(id string) <-chan evt {
	g.busMu.Lock()
	defer g.busMu.Unlock()

	if ch, ok := g.bus[id]; ok {
		return ch
	}

	ch := make(chan evt, DefaultQueueSize)
	g.bus[id] = ch
	return ch
}

func (g *Group) publish(ctx context.Context, event goevent.Event) {
	g.busMu.Lock()
	defer g.busMu.Unlock()

	for _, ch := range g.bus {
		select {
		case ch <- evt{ctx, event}:
		default:

		}
	}
}

// Close all the open channels
func (g *Group) Close() {
	for _, ch := range g.bus {
		close(ch)
	}

	g.bus = nil
}
