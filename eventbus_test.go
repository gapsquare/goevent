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
	"testing"
)

func TestEventBusError(t *testing.T) {
	var testCases = []struct {
		name              string
		err               error
		event             Event
		expectedErrorText string
	}{
		{
			"both non-nil",
			errors.New("some error"),
			NewEvent("some event type", nil),
			"some error: (some event type@0: [<nil>])",
		},
		{
			"error nil",
			nil,
			NewEvent("some event type", nil),
			"%!s(<nil>): (some event type@0: [<nil>])",
		},
		{
			"event nil",
			errors.New("some error"),
			nil,
			"some error: (%!s(<nil>))",
		},

		{
			"both nil",
			nil,
			nil,
			"%!s(<nil>): (%!s(<nil>))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			busError := EventBusError{
				Err:   tc.err,
				Event: tc.event,
			}

			if busError.Error() != tc.expectedErrorText {
				t.Errorf(
					"expected '%s', got '%s'",
					tc.expectedErrorText,
					busError.Error())
			}
		})
	}
}
