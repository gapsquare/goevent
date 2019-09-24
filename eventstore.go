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

// EventStore is an interface for event sourcing store.
type EventStore interface {
	// Save appends all events in the event stream
	Save(Events, originVersion VersionType) error
}

// EventStoreMaintainer is an interface for a maintainer of an EventStore
// NOTE: this should not be used in the real app, useful for migration tools, etc
// TODO: need more investigation
type EventStoreMaintainer interface {
	EventStore
	Replace(Event) error
	RenameEvent(from, to EventTopic) error
}
