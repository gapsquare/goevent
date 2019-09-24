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

package goevent

import (
	"context"
	"fmt"
)

// EventBusError is a async error containing the error returned from handler
// and the event that it happened on
type EventBusError struct {
	Err   error
	Event Event
}

// Error implements the Error method of the error interface.
func (e EventBusError) Error() string {
	return fmt.Sprintf("%s: (%s)", e.Err, e.Event)
}

// EventBus sends published events to one of each handler type
type EventBus interface {
	// Publish publishes the event on the bus
	Publish(context.Context, Event) error

	// AddHandler adds a handler for an event.
	AddHandler(EventMatcher, EventHandler) error //TODO: remove return error

	// Errors returns an error channel where async handling errors are sent.
	Errors() <-chan EventBusError
}
