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
	"errors"
	"fmt"
	"sync"
)

var errEventDataNotRegistered = func(topic EventTopic) error {
	return fmt.Errorf("Event data not registered for %q", topic)
}
var errRegisterEmptyTopic = errors.New("can not register empty event topic")
var errRegisterDuplicatetopic = func(topic EventTopic) error {
	return fmt.Errorf("Registering duplicate type for %q", topic)
}

var errUnregisteredNonRegistered = func(topic EventTopic) error {
	return fmt.Errorf("Unregister of non-registered topic %q", topic)
}

var eventDataFactories = make(map[EventTopic]func() EventData)
var eventDataFactoryMu sync.RWMutex

// RegisterEventData registers an event data factory for a topic. The factory is used to create event data struts
//
// An example would be:
// RegisterEventData(MyEventType, func() Event {return &MyEventData{} })
func RegisterEventData(topic EventTopic, factory func() EventData) error {
	if topic == EventTopic("") {
		return errRegisterEmptyTopic
	}

	eventDataFactoryMu.Lock()
	defer eventDataFactoryMu.Unlock()
	if _, ok := eventDataFactories[topic]; ok {
		return errRegisterDuplicatetopic(topic)
	}

	eventDataFactories[topic] = factory
	return nil
}

// UnRegisterEventData removes the registration of the event data factory for a topic.
// This is mainly useful in mainenance situations where the event data needs to be switched in a migrations.
func UnRegisterEventData(topic EventTopic) error {
	eventDataFactoryMu.Lock()
	defer eventDataFactoryMu.Unlock()

	if _, ok := eventDataFactories[topic]; !ok {
		return errUnregisteredNonRegistered(topic)
	}

	delete(eventDataFactories, topic)
	return nil
}

// CreateEventData creates an event data from a topic using the factory registered
func CreateEventData(topic EventTopic) (EventData, error) {
	eventDataFactoryMu.Lock()
	defer eventDataFactoryMu.Unlock()

	if factory, ok := eventDataFactories[topic]; ok {
		return factory(), nil
	}

	return nil, errEventDataNotRegistered(topic)
}
