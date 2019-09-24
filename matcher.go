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

// EventMatcher is a func type for matching event to a criteria
type EventMatcher func(Event) bool

// MatchAny matches any event
func MatchAny() EventMatcher {
	return func(e Event) bool {
		return true
	}
}

//MatchTopic match a specific event topic
func MatchTopic(topic EventTopic) EventMatcher {
	return func(e Event) bool {
		return e != nil && e.Topic() == topic
	}
}

// MatchAnyOf matches if any of several matchers matches
func MatchAnyOf(matches ...EventMatcher) EventMatcher {
	return func(e Event) bool {
		for _, m := range matches {
			if m(e) {
				return true
			}
		}
		return false
	}
}

// MatchAnyOfTopic matches if any of several topics matches
func MatchAnyOfTopic(topics ...EventTopic) EventMatcher {
	return func(e Event) bool {
		for _, t := range topics {
			if MatchTopic(t)(e) {
				return true
			}
		}
		return false
	}
}
