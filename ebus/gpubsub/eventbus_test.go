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
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/gapsquare/goevent/ebus"
)

// func TestGooglePubSub(t *testing.T) {
// 	// jsonKey, err := ioutil.ReadFile("pubsub-experiment.json")
// 	// if err != nil {
// 	// 	t.Fatal("Reading pubsub-key", err)
// 	// }

// 	// conf, err := google.JWTConfigFromJSON(jsonKey, pubsub.ScopePubSub, pubsub.ScopeCloudPlatform)
// 	// if err != nil {
// 	// 	t.Fatal("Create config:", err)
// 	// }

// 	if os.Getenv("PUBSUB_EMULATOR_HOST") != "" {
// 		os.Setenv("PUBSUB_EMULATOR_HOST", "")
// 		//t.Fatal("PUBSUB_EMULATOR_HOST is not set. Please setup docker pubsub first")
// 	}

// 	projectID := "experiments-gfyu"
// 	appID := "goes_testing"

// 	bus1, err := NewEventBus(projectID, appID, mocks.MockEventDataRegister, option.WithCredentialsFile("pubsub-key2.json"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	bus2, err := NewEventBus(projectID, appID, mocks.MockEventDataRegister, option.WithCredentialsFile("pubsub-key2.json"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	eventbus.AcceptanceTest(t, bus1, bus2, time.Second)
// }

func TestEventBus(t *testing.T) {
	if os.Getenv("PUBSUB_EMULATOR_HOST") == "" {
		os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8793")
		//t.Fatal("PUBSUB_EMULATOR_HOST is not set. Please setup docker pubsub first")
	}

	// Get a random app ID.
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}

	appID := "app-" + hex.EncodeToString(b)
	bus1, err := NewEventBus("project_id", appID)
	if err != nil {
		t.Fatal(err)
	}

	bus2, err := NewEventBus("project_id", appID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Run gpubsub tests")
	ebus.AcceptanceTest(t, bus1, bus2, time.Second)

}
