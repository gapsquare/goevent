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
	"reflect"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {

	event := NewEvent(TestEventTopic, &TestEventData{"event1"})
	if event.Topic() != TestEventTopic {
		t.Error("the event type should be correct:", event.Topic())
	}
	if !reflect.DeepEqual(event.Data(), &TestEventData{"event1"}) {
		t.Error("the data should be correct:", event.Data())
	}

	if event.Version() != 0 {
		t.Error("the version should be zero:", event.Version())
	}
	if event.String() != "TestEvent@0: [&{event1}]" {
		t.Error("the string representation should be correct:", event.String())
	}
}

func TestNewEventWithTimeAndVersion(t *testing.T) {

	version := VersionType(1)
	timestamp := time.Date(2019, time.January, 11, 22, 0, 0, 0, time.UTC)
	event := NewEventTimeVersion(TestEventRegisterTopic, &TestEventRegisterData{}, timestamp, version)

	if event.Timestamp() != timestamp {
		t.Errorf("timestamp differs, expected %v, actual %v", timestamp, event.Timestamp())
	}

	if event.Version() != version {
		t.Errorf("version differs, expected %v, actual %v", version, event.Version())
	}

}

func TestCreateEventData(t *testing.T) {
	data, err := CreateEventData(TestEventRegisterTopic)

	if err.Error() != errEventDataNotRegistered(TestEventRegisterTopic).Error() {
		t.Error("there should be a event not registered error:", err)
	}

	RegisterEventData(TestEventRegisterTopic, func() EventData {
		return &TestEventRegisterData{}
	})

	data, err = CreateEventData(TestEventRegisterTopic)
	if err != nil {
		t.Error("there should be no error:", err)
	}
	if _, ok := data.(*TestEventRegisterData); !ok {
		t.Errorf("the event type should be correct: %T", data)
	}

	UnRegisterEventData(TestEventRegisterTopic)
}

func TestRegisterEventEmptyName(t *testing.T) {

	err := RegisterEventData(TestEventRegisterEmptyTopic, func() EventData {
		return &TestEventRegisterEmptyData{}
	})

	if err == nil {
		t.Errorf("goevent: registering empty topic not allowed, excepted errRegisterEmptyTopic error")
	}

	if err != errRegisterEmptyTopic {
		t.Errorf("goevent: expected errRegisterEmptyTopic but got %v", err)
	}
}

func TestRegisterEventTwice(t *testing.T) {

	RegisterEventData(TestEventRegisterTwiceTopic, func() EventData {
		return &TestEventRegisterTwiceData{}
	})
	err := RegisterEventData(TestEventRegisterTwiceTopic, func() EventData {
		return &TestEventRegisterTwiceData{}
	})

	if err == nil {
		t.Errorf("goevent: registering topic twice not allowed, excepted errRegisterDuplicatetopic error")
	}

	if err.Error() != errRegisterDuplicatetopic(TestEventRegisterTwiceTopic).Error() {
		t.Errorf("goevent: expected errRegisterDuplicatetopic but got %v", err)
	}
}

func TestUnregisterEventEmptyName(t *testing.T) {

	err := UnRegisterEventData(TestEventUnregisterEmptyTopic)

	if err == nil {
		t.Errorf("goevent: registering topic twice not allowed, excepted errUnregisteredNonRegistered error")
	}

	if err.Error() != errUnregisteredNonRegistered(TestEventUnregisterEmptyTopic).Error() {
		t.Errorf("goevent: expected errUnregisteredNonRegistered but got %v", err)
	}
}

func TestUnregisterEventTwice(t *testing.T) {

	RegisterEventData(TestEventUnregisterTwiceTopic, func() EventData {
		return &TestEventUnregisterTwiceData{}
	})
	UnRegisterEventData(TestEventUnregisterTwiceTopic)
	err := UnRegisterEventData(TestEventUnregisterTwiceTopic)
	if err == nil {
		t.Errorf("goevent: registering topic twice not allowed, excepted errUnregisteredNonRegistered error")
	}

	if err.Error() != errUnregisteredNonRegistered(TestEventUnregisterTwiceTopic).Error() {
		t.Errorf("goevent: expected errUnregisteredNonRegistered but got %v", err)
	}
}

const (
	TestEventTopic                EventTopic = "TestEvent"
	TestEventRegisterTopic        EventTopic = "TestEventRegister"
	TestEventRegisterEmptyTopic   EventTopic = ""
	TestEventRegisterTwiceTopic   EventTopic = "TestEventRegisterTwice"
	TestEventUnregisterEmptyTopic EventTopic = ""
	TestEventUnregisterTwiceTopic EventTopic = "TestEventUnregisterTwice"
)

type TestEventData struct {
	Content string
}

type TestEventRegisterData struct{}

type TestEventRegisterEmptyData struct{}

type TestEventRegisterTwiceData struct{}

type TestEventUnregisterTwiceData struct{}
