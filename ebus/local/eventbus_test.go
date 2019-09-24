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
	"testing"
	"time"

	"github.com/isuruceanu/goevent/ebus"
)

func TestEventBus(t *testing.T) {
	group := NewGroup()
	if group == nil {
		t.Fatal("there should be a group")
	}

	bus1, err := NewEventBus(group)
	if err != nil {
		t.Error(err)
	}

	bus2, err := NewEventBus(group)
	if err != nil {
		t.Error(err)
	}

	t.Log("Run local test")

	ebus.AcceptanceTest(t, bus1, bus2, time.Second)
	bus1.Close()
	bus2.Close()
	bus1.Wait()
	bus2.Wait()
}
