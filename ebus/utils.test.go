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
	"fmt"
	"reflect"

	"github.com/gapsquare/goevent"
)

// CompareEvents copares two events, ignoring their version and timestamp
func CompareEvents(e1, e2 goevent.Event) error {
	if e1.Topic() != e2.Topic() {
		return fmt.Errorf("incorrect event topic: %s (should be %s)", e1.Topic(), e2.Topic())
	}

	if !reflect.DeepEqual(e1.Data(), e2.Data()) {
		return fmt.Errorf("incorrect event data: %s (should be %s", e1.Data(), e2.Data())
	}
	return nil
}

// EqualEvents compares two slice of events
func EqualEvents(evts1, evts2 []goevent.Event) bool {
	if len(evts1) != len(evts2) {
		return false
	}

	for i, e1 := range evts1 {
		e2 := evts2[i]

		if e1.Topic() != e2.Topic() {
			return false
		}

		if !reflect.DeepEqual(e1.Data(), e2.Data()) {
			return false
		}

		if e1.Timestamp() != e2.Timestamp() {
			return false
		}

		if e1.Version() != e2.Version() {
			return false
		}
	}

	return true
}
