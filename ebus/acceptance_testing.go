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

package ebus

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/isuruceanu/goevent"

	"github.com/kr/pretty"
)

// AcceptanceTest is the acceptance test that all implementations of EventBus
// should pass. it should manually be called from a test case in each implementation
//
// func TestEventBus(t *testing.T) {
// 	bus1 := NewEventBus()
// 	bus2 := NewEventBus()
// 	eventbus.AcceptanceTest(t, bus1, bus2)
// }
func AcceptanceTest(t *testing.T, bus1, bus2 goevent.EventBus, timeout time.Duration) {

	// Panic on nil matcher
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("adding a nil matcher should panic")
			}
		}()
		bus1.AddHandler(nil, NewEventHandler("panic"))

	}()

	// Panic on nil handler
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("adding a nil handler should panic")
			}
		}()
		bus1.AddHandler(goevent.MatchAny(), nil)
	}()

	// Panic on multiple registrations
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("adding multiple handlers should panic")
			}
		}()
		bus1.AddHandler(goevent.MatchAny(), NewEventHandler("multi"))
		bus1.AddHandler(goevent.MatchAny(), NewEventHandler("multi"))
	}()

	//ctx := mocks.WithContextOne(context.Background(), "testval")
	ctx := context.Background()

	// Without handler
	//id, _ := uuid.Parse("b2238a5b-c5bb-5aa0-3acb-1a5b6b8d1234")
	timestamp := time.Date(2019, time.January, 11, 22, 0, 0, 0, time.UTC)
	event1 := goevent.NewEventForAggregate(Topic, &EventData{Content: "event1"}, timestamp, goevent.VersionType(1))
	if err := bus1.Publish(ctx, event1); err != nil {
		t.Error("there should be no error: ", err)
	}

	const (
		handlerName = "handler"
	)

	handlerBus1 := NewEventHandler(handlerName)
	handlerBus2 := NewEventHandler(handlerName)
	anotherHandlerBus2 := NewEventHandler("another_handler")
	bus1.AddHandler(goevent.MatchAny(), handlerBus1)
	bus2.AddHandler(goevent.MatchAny(), handlerBus2)
	bus2.AddHandler(goevent.MatchAny(), anotherHandlerBus2)

	if err := bus1.Publish(ctx, event1); err != nil {
		t.Error("there should be no error:", err)
	}

	// Check for correct event in handler 1 or 2.
	expectedEvents := []goevent.Event{event1}
	if !(handlerBus1.Wait(timeout) || handlerBus2.Wait(timeout)) {
		t.Error("did not receive event in time")
	}

	if !(EqualEvents(handlerBus1.Events, expectedEvents) ||
		EqualEvents(handlerBus2.Events, expectedEvents)) {
		t.Error("the events were incorrect:")
		t.Log(handlerBus1.Events)
		t.Log(handlerBus2.Events)
		if len(handlerBus1.Events) == 1 {
			t.Log(pretty.Sprint(handlerBus1.Events[0]))
		}
		if len(handlerBus2.Events) == 1 {
			t.Log(pretty.Sprint(handlerBus2.Events[0]))
		}
	}
	if EqualEvents(handlerBus1.Events, handlerBus2.Events) {
		t.Error("only one handler should receive the events")
	}

	// Check the other handler.
	if !anotherHandlerBus2.Wait(timeout) {
		t.Error("did not receive event in time")
	}
	if !EqualEvents(anotherHandlerBus2.Events, expectedEvents) {
		t.Error("the events were incorrect:")
		t.Log(anotherHandlerBus2.Events)
	}

	// Test async errors from handlers.
	errorHandler := NewEventHandler("error_handler")
	errorHandler.Err = errors.New("handler error")
	bus1.AddHandler(goevent.MatchAny(), errorHandler)
	if err := bus1.Publish(ctx, event1); err != nil {
		t.Error("there should be no error:", err)
	}
	select {
	case <-time.After(time.Second):
		t.Error("there should be an async error")
	case err := <-bus1.Errors():
		// Good case.
		if err.Error() != "could not handle event (error_handler): handler error: (Event@1: [&{event1}])" {
			t.Error(err, "wrong error sent on event bus")
		}
	}

}
