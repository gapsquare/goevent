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
	"fmt"
	"time"
)

// EventTopic is the topic of an event, used as its unique identifier
type EventTopic string

// EventData is any additional data for an event
type EventData interface{}

// TimeNow mockable version of time.Now
var TimeNow = time.Now

// Event is a domain event describing changes happened
type Event interface {
	// Topic returns the event topic
	Topic() EventTopic

	// Data returns event data
	Data() EventData

	// Timestamp of the event was created
	Timestamp() time.Time

	// Version of the aggregate for this event (after it was applied)
	Version() VersionType

	// String stringer implementation
	String() string
}

// Events list of Event s
type Events []Event

// NewEvent creates a new event for topic with data
func NewEvent(topic EventTopic, data EventData) Event {
	return event{
		topic:     topic,
		data:      data,
		timestamp: TimeNow(),
	}
}

// NewEventTimeVersion creates a new event for an aggregate object
func NewEventTimeVersion(topic EventTopic, data EventData, timestamp time.Time, version VersionType) Event {
	return event{
		topic:     topic,
		data:      data,
		timestamp: timestamp,
		version:   version,
	}
}

type event struct {
	topic     EventTopic
	data      EventData
	timestamp time.Time
	version   VersionType
}

func (e event) Topic() EventTopic    { return e.topic }
func (e event) Data() EventData      { return e.data }
func (e event) Timestamp() time.Time { return e.timestamp }
func (e event) Version() VersionType { return e.version }

func (e event) String() string {
	return fmt.Sprintf("%s@%d: [%v]", e.topic, e.version, e.data)
}
