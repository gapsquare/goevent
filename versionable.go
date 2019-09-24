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

import "errors"

// ErrVersionMismatched error for version mismatched
var ErrVersionMismatched = errors.New("Version mismatched")

// Versionable is an item that has a version number
type Versionable interface {
	Version() VersionType
}

// VersionType version type
type VersionType uint64
