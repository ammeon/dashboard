// Copyright 2015 Google Inc. All Rights Reserved.
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

package common

import (
	"time"
)

type Release struct {
	Name      string    `json:"name"`
	Time      time.Time `json:"time"`
	Namespace string    `json:"namespace"`
	Status    string    `json:"status"`
}

type ReleaseList struct {
	// Items is the list of deployments.
	Items []Release `json:"items"` // TODO: Releases
}